package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return &RedisClient{client: client}
}
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rbd := NewRedisClient()
	if err := rbd.Ping(ctx); err != nil {
		log.Fatalln(err)
	}
	log.Println("Redis connected")
	r := gin.Default()
	r.Use(LoggerMiddleware())
	r.POST("/files", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = c.SaveUploadedFile(file, "uploads/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save file",
			})
			return
		}
		err = rbd.client.XAdd(c.Request.Context(), &redis.XAddArgs{
			Stream: "file-events",
			Values: map[string]interface{}{
				"filename": file.Filename,
			},
		}).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to publish event",
			})
			return
		}

		log.Println("File name:", file.Filename)
		log.Println("File upload request received")
		c.String(http.StatusOK, "File uploaded")
	})
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Request:", c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}
