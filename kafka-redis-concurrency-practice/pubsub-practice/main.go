package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)

	}
	log.Println("Redis Connected")

	pubsub := redisClient.Subscribe(ctx, "notifications")
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	received := 0
	go func() {
		defer wg.Done()
		ch := pubsub.Channel()
		for {
			select {

			case msg := <-ch:
				log.Printf("received: %s", msg.Payload)

				received++

				if received == 6 {
					cancel()
				}

			case <-ctx.Done():
				log.Println("subscriber stopped")
				return

			}
		}
	}()
	messages := []string{"msg-1", "msg-2", "msg-3", "message-1",
		"message-2",
		"message-3"}
	for _, message := range messages {
		err := redisClient.Publish(ctx, "notifications", message)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("published: %s", message)

	}

	wg.Wait()
}
