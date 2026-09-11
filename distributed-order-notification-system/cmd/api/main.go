package main

import (
	"context"
	"distributed-order-notification-system/internal/config"
	"distributed-order-notification-system/internal/handler"
	"distributed-order-notification-system/internal/middleware"
	"distributed-order-notification-system/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("PostgreSQL connected successfully")
	redisclient, err := config.ConnectRedis(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	err = redisclient.Set(context.Background(), "test1", "hello", 0).Err()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data stored in Redis")
	data, err := redisclient.Get(context.Background(), "test1").Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data got in Redis", data)
	service.RunWorkers()
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/orders", handler.CreateOrder)
	log.Println("Server running on :8080")
	handler := http.DefaultServeMux
	wrappedHandler := middleware.Logger(handler)
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}
