package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/work", func(c *gin.Context) {
		lockValue := uuid.NewString()
		locked, err := redisClient.AcquireLock(
			c.Request.Context(),
			"work-lock",
			lockValue,
			10*time.Second,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed to acquire lock",
			})
			return
		}

		if !locked {
			c.JSON(409, gin.H{
				"error": "another request is processing",
			})
			return
		}
		time.Sleep(5 * time.Second)

		if err := redisClient.ReleaseLock(c.Request.Context(), "work-lock", lockValue); err != nil {
			c.JSON(500, gin.H{
				"error": "failed to release lock",
			})
			return
		}

		c.JSON(200, gin.H{"message": "work completed"})
	})
	r.Run(":8081")

}
