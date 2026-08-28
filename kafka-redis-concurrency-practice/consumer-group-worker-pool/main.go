package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"log"
	"sync"
)

func main() {
	ctx := context.Background()
	writer := kafka.NewWriter("worker-pool-practice")
	defer writer.Close()

	jobs := make(chan string)
	workerPool := NewWorkerPool(3)
	var wg sync.WaitGroup
	wg.Add(1)

	go workerPool.Start(jobs, &wg)
	reader := kafka.NewReader("worker-pool-practice",
		"worker-pool-group")
	defer reader.Close()
	for i := 1; i <= 10; i++ {
		message := []byte(fmt.Sprintf("job-%d", i))
		err := writer.WriteMessage(ctx, message)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("produced:%s", message)
	}

	for i := 0; i < 10; i++ {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("consumed: %s", msg.Value)
		jobs <- string(msg.Value)
		if err := reader.CommitMessage(ctx, msg); err != nil {
			log.Fatal(err)
		}

	}
	close(jobs)
	wg.Wait()

}
