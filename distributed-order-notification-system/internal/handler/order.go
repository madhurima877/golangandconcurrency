package handler

import (
	"distributed-order-notification-system/internal/service"
	"encoding/json"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var order struct {
		Product string `json:"product"`
		Amount  int    `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	go service.ProcessOrder(1)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order created",
		"product": order.Product,
		"amount":  order.Amount,
	})
}
