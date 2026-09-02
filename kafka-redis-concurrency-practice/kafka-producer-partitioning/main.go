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
	writer := kafka.NewKeyWriter("rebalance-practice")
	defer writer.Close()
	log.Println("producer started")

	var wg sync.WaitGroup
	go kafka.ReaderWithWorker(ctx, "rebalance-practice", "partitioning-group", "worker-1")

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			err := writer.WriteMessageWithKey(
				ctx,
				[]byte(fmt.Sprintf("user-%d", i)),
				[]byte(fmt.Sprintf("job-%d", i)),
			)

			if err != nil {
				log.Printf("failed job-%d: %v", i, err)
				return
			}

			log.Printf("produced job-%d", i)
		}(i)
	}

	wg.Wait()

	select {}
}
