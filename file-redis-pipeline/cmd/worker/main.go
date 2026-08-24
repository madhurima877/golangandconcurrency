package main

import (
	"context"
	"log"
	"strings"
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
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rbd := NewRedisClient()
	if err := rbd.Ping(ctx); err != nil {
		log.Fatalln(err)
	}
	log.Println("Worker connected")
	err := rbd.client.XGroupCreateMkStream(ctx, "file-events", "file-workers", "0").Err()

	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Fatalln(err)
	}
	log.Println("Worker connected")

	for {
		streams, err := rbd.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "file-workers",
			Consumer: "worker-1",
			Streams:  []string{"file-events", ">"},
			Block:    0,
		}).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				filename := message.Values["filename"]

				log.Println("Processing file:", filename)

				time.Sleep(3 * time.Second)

				log.Println("File processed:", filename)

				err := rbd.client.XAck(
					ctx,
					"file-events",
					"file-workers",
					message.ID,
				).Err()

				if err != nil {
					log.Println("Failed to acknowledge:", err)
					continue
				}

				log.Println("Message acknowledged:", message.ID)
			}
		}
	}
	// log.Println(streams)
}
