package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var processed = make(map[string]bool)

var attempts int

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if processed[key] {
		fmt.Fprintln(w, "already pricessed")
		return
	}
	attempts++
	if attempts <= 2 {
		http.Error(w, "temporary failure", http.StatusInternalServerError)
		return
	}
	processed[key] = true
	fmt.Fprintln(w, "order created")
}
func main() {
	http.HandleFunc("/orders", handler)
	go func() {
		log.Println("server started on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}

	}()
	time.Sleep(500 * time.Millisecond)
	log.Println("sending same request again")

	log.Println("first request")
	err := doRequest()
	if err != nil {
		log.Println("request failed:", err)
	}

	log.Println("second request with same idempotency key")
	err = doRequest()
	if err != nil {
		log.Println("request failed:", err)
	}

}
func doRequest() error {
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(
			http.MethodGet,
			"http://localhost:8080/orders",
			nil,
		)
		if err != nil {
			return err
		}

		req.Header.Set("Idempotency-Key", "order-123")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}
			log.Printf("attempt %d succeeded: %s", attempt, body)
			return nil
		}

		resp.Body.Close()

		log.Printf("attempt %d failed", attempt)

		if attempt < 3 {
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("request failed after 3 attempts")
}
