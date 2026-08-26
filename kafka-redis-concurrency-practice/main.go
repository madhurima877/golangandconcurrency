package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"golangandconcurrency/kafka-redis-concurrency-practice/worker"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisClient := redis.New()

	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	stream := "jobs"
	group := "workers"
	err := redisClient.CreateConsumerGroup(ctx, stream, group)
	if err != nil {
		log.Println("consumer group:", err)
	}
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		w := worker.New(redisClient, fmt.Sprintf("worker-%d", i))
		go func() {
			defer wg.Done()
			w.Start(ctx, stream, group)
		}()
	}
	time.Sleep(time.Second)

	for i := 1; i <= 10; i++ {
		err := redisClient.AddMessage(ctx, stream, map[string]interface{}{
			"job": fmt.Sprintf("job-%d", i),
		})
		if err != nil {
			log.Println("add message error:", err)
		}
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutdown signal received")

	cancel()
	wg.Wait()

	fmt.Println("all workers stopped")

}

func main2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("cancelling context...")
		cancel()
	}()
	writer := kafka.NewWriter("practice-topic")

	err := writer.WriteMessage(ctx, []byte("hello kafka"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("message produced")
	reader := kafka.NewReader("practice-topic", "practice-group")
	msg, err := reader.FetchMessage(ctx)
	if err != nil {
		if err == context.Canceled {
			fmt.Println("consumer stopped because context was cancelled")
			return
		}

		log.Fatal(err)
	}
	log.Printf("received: %s\n", string(msg.Value))
	fmt.Println("processing...")
	time.Sleep(2 * time.Second)
	err = reader.CommitMessage(ctx, msg)
	if err != nil {
		if err == context.Canceled {
			fmt.Println("commit stopped because context was cancelled")
			return
		}

		log.Fatal(err)
	}

	fmt.Println("Offset Committed")

	fmt.Println("processing finished — NOT committing offset")
}
