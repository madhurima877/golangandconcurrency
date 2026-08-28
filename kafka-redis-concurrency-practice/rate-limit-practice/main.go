package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.Use(RateLimit(redisClient))
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "request allowed"})
	})
	r.Run(":8082")
}
