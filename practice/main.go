package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// func main() {
// 	sem := make(chan struct{}, 3)
// 	var wg sync.WaitGroup
// 	for i := 0; i <= 10; i++ {
// 		wg.Add(1)
// 		sem <- struct{}{}

// 		go func() {
// 			defer wg.Done()
// 			defer func() {
// 				<-sem
// 			}()
// 			fmt.Println("Running:", i)
// 			time.Sleep(time.Second)
// 			fmt.Println("Finished:", i)

// 		}()
// 	}
// 	wg.Wait()
// }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	jobs := make(chan int)
	results := make(chan int)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg, results, jobs, ctx)
	}
	go func() {
		defer close(jobs)
		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- i:
			}

		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Shutdown signal received")
		cancel()

	}()

	go func() {
		wg.Wait()
		close(results)
	}()
	for result := range results {
		fmt.Println("Result:", result)
	}

	log.Println("Server shutdown gracefully")

}
func worker(id int, wg *sync.WaitGroup, results chan<- int, jobs <-chan int, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				log.Printf("Worker %d: jobs closed\n", id)
				return
			}
			log.Printf("Worker %d processing job %d\n", id, job)
			time.Sleep(500 * time.Millisecond)
			result := job * job
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}

		}
	}

}
