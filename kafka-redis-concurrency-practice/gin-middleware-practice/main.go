package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(RequestID())
	r.Use(StructuredLogger())
	r.GET("/hello", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "request ID missing",
			})
			return
		}
		ProcessRequest(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{
			"message":    "hello",
			"request_id": requestID,
		})
	})
	r.Run(":8083")
}

func ProcessRequest(ctx context.Context) {
	requestID := ctx.Value("request_id")
	log.Printf(
		"processing request_id=%v",
		requestID,
	)
}
