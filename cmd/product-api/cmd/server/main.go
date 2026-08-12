package main

import (
	"context"
	"log"
	"product-api/internal/handler"
	"product-api/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := repository.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := gin.Default()
	cache := repository.NewRedisCache()
	if err := cache.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	repo := repository.NewProductRepository(db)
	handler := handler.NewProductHandler(repo, cache)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/products", handler.Create)
	r.GET("/products/:id", handler.Get)
	log.Println("server is running on :8080")
	r.Run(":8080")
}
