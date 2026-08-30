package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Redis Connected")
	id, err := redisClient.AddToStream(ctx, "orders", map[string]interface{}{
		"order_id": "order-1",
		"status":   "created",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("message added: %s", id)

	err = redisClient.CreateConsumerGroup(ctx, "orders", "order-workers")
	err = redisClient.CreateConsumerGroup(
		ctx,
		"orders",
		"order-workers",
	)

	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatal(err)
	}

	log.Println("consumer group ready")
	log.Println("consumer group created")

	messages, err := redisClient.ReadGroup(ctx, "orders", "order-workers", "worker-1")
	if err != nil {
		log.Fatal(err)
	}
	for _, message := range messages {
		log.Printf("received: %v", message)
	}
	err = redisClient.Ack(
		ctx,
		"orders",
		"order-workers",
		"1788070138476-0",
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("message acknowledged")
	err = redisClient.Claim(
		ctx,
		"orders",
		"order-workers",
		"worker-2",
		time.Second,
		"1788072892036-0",
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("message claimed by worker-2")
}
