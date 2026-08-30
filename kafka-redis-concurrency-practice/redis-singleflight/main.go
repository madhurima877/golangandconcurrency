package main

import (
	"context"
	"log"
	"time"

	"golangandconcurrency/kafka-redis-concurrency-practice/redis"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

func main() {
	ctx := context.Background()

	redisClient := redis.New()

	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis Connected")

	router := gin.Default()
	router.GET("/user/:id", func(c *gin.Context) {
		ctx := c.Request.Context()
		userID := c.Param("id")
		key := "user:" + userID
		data, err := redisClient.Get(ctx, key)
		if err == nil {
			c.JSON(200, gin.H{"data": data})
			return
		}
		result, err, _ := group.Do(key, func() (interface{}, error) {
			data, err := fetchFromDB(ctx)

			if err != nil {
				return nil, err
			}
			err = redisClient.Set(ctx, "user:1", data, 30*time.Second)
			if err != nil {
				return nil, err
			}
			return data, nil
		})
		if err != nil {
			log.Printf("failed: %v", err)
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}

		data = result.(string)

		c.JSON(200, gin.H{"data": data})
	})
	router.Run(":8084")
}
func fetchFromDB(ctx context.Context) (string, error) {
	log.Println("DB CALL STARTED")

	select {
	case <-time.After(2 * time.Second):
		log.Println("DB CALL FINISHED")
		return "user-data", nil

	case <-ctx.Done():
		return "", ctx.Err()
	}
}
