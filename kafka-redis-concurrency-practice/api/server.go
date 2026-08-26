package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/api/fileworker"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

type Server struct {
	redisClient *redis.Client
}

func NewServer(redisClient *redis.Client) *Server {
	return &Server{redisClient: redisClient}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	server := NewServer(redisClient)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)

		worker := fileworker.New(
			redisClient,
			fmt.Sprintf("worker-%d", i),
		)

		go func(w *fileworker.Worker) {
			defer wg.Done()
			w.Start(ctx)
		}(worker)
	}

	r := gin.Default()
	r.POST("/upload/file", server.Upload)

	go r.Run(":8080")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-shutdown

	fmt.Println("shutdown signal received")

	cancel()

	wg.Wait()

	fmt.Println("all workers stopped")

}
func (s *Server) Upload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}
	err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to save file",
		})
		return
	}

	err = s.redisClient.PushJob(c.Request.Context(), "file-jobs", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to queue file",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "file uploaded and queued",
		"file":    file.Filename,
	})
}
