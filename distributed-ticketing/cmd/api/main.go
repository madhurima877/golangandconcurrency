package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	fmt.Println("API running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
