package main

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var produced int64
	var processed int64
	var dropped int64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(time.Second)
		log.Println("cancelling all workers")
		cancel()
	}()
	jobs := make(chan int, 5)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go Consumer(i, jobs, &wg, ctx, &processed)
	}
	go Producer(jobs, ctx, &produced, &dropped)
	wg.Wait()
	log.Println("----- Metrics -----")
	log.Printf("produced: %d", atomic.LoadInt64(&produced))
	log.Printf("processed: %d", atomic.LoadInt64(&processed))
	log.Printf("dropped: %d", atomic.LoadInt64(&dropped))
}
func Producer(jobs chan<- int, ctx context.Context, produced, dropped *int64) {
	defer close(jobs)

	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			log.Println("producer stopped")
			return

		case jobs <- i:
			atomic.AddInt64(produced, 1)
			log.Printf("produced job %d", i)

		case <-time.After(500 * time.Millisecond):
			atomic.AddInt64(dropped, 1)
			log.Printf("dropping job %d: queue still full", i)
		}
	}
}
func Consumer(i int, jobs <-chan int, wg *sync.WaitGroup, ctx context.Context, processed *int64) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopped", i)
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				log.Printf("worker %d stopped", i)
				return
			}
			atomic.AddInt64(processed, 1)

			log.Println("Job ID ", job, " done by worker id ", i)
		}
	}
}
