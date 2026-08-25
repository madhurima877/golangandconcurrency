package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func processFile(filename string) error {
	log.Println("Processing file:", filename)

	time.Sleep(3 * time.Second)

	log.Println("File processed successfully:", filename)

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	rbd := NewRedisClient()

	if err := rbd.Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("Worker connected")

	err := rbd.client.XGroupCreateMkStream(
		ctx,
		"file-events",
		"file-workers",
		"0",
	).Err()

	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Fatalln(err)
	}

	consumerName := os.Getenv("CONSUMER_NAME")

	if consumerName == "" {
		consumerName = "worker-1"
	}

	for {
		streams, err := rbd.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "file-workers",
			Consumer: consumerName,
			Streams:  []string{"file-events", ">"},
			Block:    0,
		}).Result()

		if err != nil {
			if ctx.Err() != nil {
				log.Println("Worker shutting down")
				return
			}

			log.Println(err)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {

				filename, ok := message.Values["filename"].(string)

				if !ok {
					log.Println("Invalid filename")
					continue
				}

				maxRetries := 3
				success := false

				for attempt := 1; attempt <= maxRetries; attempt++ {

					log.Printf(
						"Attempt %d for file: %s",
						attempt,
						filename,
					)

					err := processFile(filename)

					if err == nil {
						success = true
						break
					}

					log.Println("File processing failed:", err)

					if attempt < maxRetries {
						time.Sleep(2 * time.Second)
					}
				}

				if !success {
					_, err := rbd.client.XAdd(ctx, &redis.XAddArgs{
						Stream: "file-events-dlq",
						Values: map[string]interface{}{
							"filename": filename,
						},
					}).Result()

					if err != nil {
						log.Println("Failed to send message to DLQ:", err)
						continue
					}

					log.Println("Message moved to DLQ")
				}

				// TEMPORARY TEST:
				// Stop before XACK so the message stays pending.
				log.Println("Stopping before acknowledgement for testing")
				return
			}
		}
	}
}
