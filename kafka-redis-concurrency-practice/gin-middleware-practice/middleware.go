package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		ctx := context.WithValue(
			c.Request.Context(),
			requestIDKey,
			requestID,
		)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)
		c.Next()

	}
}

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get(requestIDKey)
		log.Printf(
			"request_id=%s method=%s path=%s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
		)
		c.Next()
	}
}
