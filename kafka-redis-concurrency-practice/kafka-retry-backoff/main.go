package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writer := kafka.NewWriter("retry-practice")
	defer writer.Close()

	err := writer.WriteMessageWithKey(ctx, []byte("job-1"), []byte("process-me"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("message produced")

	reader := kafka.NewReader("retry-practice", "retry-practice-group")
	defer reader.Close()
	msg, err := reader.FetchMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		time.Sleep(1500 * time.Millisecond)
		cancel()
	}()
	log.Printf(
		"consumed: key=%s value=%s",
		msg.Key,
		msg.Value,
	)
	err = retry(ctx, 3)

	if err != nil {
		log.Printf("processing failed: %v", err)
	}
}
func processMessage() error {
	return fmt.Errorf("temporary processing failure")
}
func retry(ctx context.Context, maxRetries int) error {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := processMessage()
		if err == nil {
			return nil
		}
		log.Printf("attempt %d failed: %v", attempt, err)
		if attempt == maxRetries {
			return err
		}
		delay := time.Duration(1<<uint(attempt-1)) * time.Second
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			log.Println("retrying....")
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
	return nil
}
