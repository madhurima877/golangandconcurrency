package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"sync"
)

func main() {
	ctx := context.Background()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Redis Connected")
	queue := "jobs"
	for i := 1; i <= 10; i++ {
		job := fmt.Sprintf("job=%d", i)
		if err := redisClient.Push(ctx, queue, job); err != nil {
			log.Fatal(err)
		}
		log.Printf("pushed: %s", job)
	}
	jobs := make(chan string, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	workerPool := NewWorkerPool(3)
	go workerPool.Start(jobs, &wg)
	for i := 1; i <= 10; i++ {
		job, err := redisClient.Pop(ctx, queue)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("popped: %s", job)
		log.Printf("trying to enqueue: %s", job)

		jobs <- job

		log.Printf("enqueued: %s", job)
	}
	close(jobs)
	wg.Wait()
}
