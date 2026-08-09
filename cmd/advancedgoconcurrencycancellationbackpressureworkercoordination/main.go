package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID       int
	Priority int
}

func worker(
	ctx context.Context,
	id int,
	jobs <-chan Job,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			fmt.Printf("Worker %d processing Job %d\n", id, job.ID)

			// Simulate work
			select {
			case <-time.After(time.Duration(job.Priority) * 200 * time.Millisecond):
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jobs := make(chan Job, 5)

	var wg sync.WaitGroup

	// Limited concurrency
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg)
	}

	// Producer
	for i := 1; i <= 20; i++ {
		select {
		case jobs <- Job{
			ID:       i,
			Priority: (i % 3) + 1,
		}:
			fmt.Println("Queued:", i)

		case <-ctx.Done():
			fmt.Println("Producer cancelled")
			close(jobs)
			wg.Wait()
			return
		}
	}

	close(jobs)
	wg.Wait()

	fmt.Println("All workers stopped")
}