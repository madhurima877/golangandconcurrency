package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Redisclient struct {
	cleint *redis.Client
}

func NewClient() *Redisclient {
	client := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	return &Redisclient{
		cleint: client,
	}
}
func (r *Redisclient) Ping(ctx context.Context) error {
	return r.cleint.Ping(ctx).Err()
}
func main() {
	ctx := context.Background()

	rdb := NewClient()
	if err := rdb.Ping(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Redis connected")
	http.HandleFunc("/orders/process", rdb.ProcessOrderHandler)
	http.HandleFunc("/orders/status", rdb.GetOrderStatusHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
func (rdb *Redisclient) ProcessOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	orderId := r.URL.Query().Get("order_id")
	if orderId == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}
	log.Println("Processing order")
	token := uuid.NewString()
	lockKey := "order:lock:" + orderId
	claimed, err := rdb.cleint.SetNX(r.Context(), lockKey, token, 30*time.Second).Result()
	if err != nil {
		log.Printf("redis error for %s: %v", orderId, err)
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	if !claimed {

		http.Error(w, "Order is already being processed", http.StatusConflict)
		return
	}

	statusKey := "order:status:" + orderId

	if err := rdb.cleint.Set(
		r.Context(),
		statusKey,
		"processing",
		0,
	).Err(); err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	defer func() {
		value, err := rdb.cleint.Get(r.Context(), lockKey).Result()
		if err != nil {
			return
		}
		if value == token {
			rdb.cleint.Del(r.Context(), lockKey)
		}
	}()
	log.Printf("Lock acquired for order: %s", orderId)
	time.Sleep(5 * time.Second)
	if err := rdb.cleint.Set(r.Context(), statusKey, "completed", 0).Err(); err != nil {
		http.Error(w, "failed to update order status", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Order processed")
}
func (rdb *Redisclient) GetOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderId := r.URL.Query().Get("order_id")
	if orderId == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}
	statusKey := "order:status:" + orderId
	status, err := rdb.cleint.Get(r.Context(), statusKey).Result()
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Order %s status: %s\n", orderId, status)
}
