package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signchan := make(chan os.Signal, 1)
	signal.Notify(signchan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signchan)
	go func() {
		<-signchan
		log.Println("shutdown signal received")
		cancel()
	}()

	reader := kafka.NewReader(
		"rebalance-practice",
		"rebalance-group",
	)
	defer reader.Close()
	log.Println("consumer started")

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Consumer shutting down")
				break
			}
			log.Printf("consumer error: %v", err)
			continue
		}
		log.Printf(
			"received: partition=%d offset=%d value=%s",
			msg.Partition,
			msg.Offset,
			msg.Value,
		)
		if err := reader.CommitMessage(ctx, msg); err != nil {
			log.Fatal(err)
		}

	}
}
