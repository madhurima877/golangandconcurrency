package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"log"
	"sync"
)

func main() {
	ctx := context.Background()
	writer := kafka.NewWriter("partition-practice")
	defer writer.Close()
	messages := []struct {
		key   string
		value string
	}{
		{"user-1", "message from user-1"},
		{"user-2", "message from user-2"},
		{"user-3", "message from user-3"},
	}
	var wg sync.WaitGroup
	for _, message := range messages {
		wg.Add(1)
		go func(msg struct {
			key   string
			value string
		}) {
			defer wg.Done()
			err := writer.WriteMessageWithKey(
				ctx,
				[]byte(msg.key),
				[]byte(msg.value),
			)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("produced key=%s\n", msg.key)

		}(message)

	}

	wg.Wait()
	Consume()
}
