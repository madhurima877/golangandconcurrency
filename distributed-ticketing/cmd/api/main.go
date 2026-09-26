package main

import (
	"distributed-ticketing/internal/db"
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	fmt.Println("API running on :8080")
	if err := http.ListenAndServe(":8080"); err != nil {
		panic(err)
	}
}
