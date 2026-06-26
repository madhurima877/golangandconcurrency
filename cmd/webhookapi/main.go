package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookPayload struct {
	Event string `json:"event"`
	ID    string `json:"id"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Println("Event:", payload.Event)
	fmt.Println("ID:", payload.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

//add these
// Authentication/signature verification
// Retry handling
// Idempotency
// Asynchronous processing (Kafka, queues)
// What if the sender sends it twice?
// What if someone fakes the request?
// What if your server is down?
