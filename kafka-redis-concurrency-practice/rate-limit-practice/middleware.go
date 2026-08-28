package main

import (
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate-limit:" + c.ClientIP()
		ctx := c.Request.Context()
		count, err := redisClient.Increment(ctx, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limit check failed",
			})
			c.Abort()
			return
		}
		if count == 1 {
			err := redisClient.Expire(ctx, key, 10*time.Second)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to set rate limit expiry",
				})
				c.Abort()
				return
			}
		}

		if count > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
