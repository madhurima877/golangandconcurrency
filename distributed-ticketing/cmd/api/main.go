package main

import (
	"distributed-ticketing/internal/db"
	"distributed-ticketing/internal/handler"
	"distributed-ticketing/internal/repository"
	"distributed-ticketing/internal/service"
	"distributed-ticketing/internal/worker"
	"fmt"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	database, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer database.Close()
	outboxWorker := worker.NewOutboxWorker(database)
	if err := outboxWorker.Process(); err != nil {
		log.Println(err)
	}

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/users", userHandler.CreateUser)

	fmt.Println("API running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
