package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Job struct {
	ID int
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	jobQueue := make(chan Job, 10)
	log.Println("job queue created", cap(jobQueue))
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobQueue, &wg)
	}
	for i := 1; i <= 10; i++ {
		jobQueue <- Job{ID: i}
	}
	log.Println("jobs submitted")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	log.Println("shutdown signal received")

	// cancel()
	close(jobQueue)

	wg.Wait()

	log.Println("worker finished")
}

// func worker(id int, ctx context.Context, jobQueue <-chan Job, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case job, ok := <-jobQueue:
// 			if !ok {
// 				log.Printf("worker-%d stopped", id)
// 				return
// 			}

// 			log.Printf("worker-%d processing job-%d", id, job.ID)
// 			time.Sleep(1500 * time.Millisecond)
// 		case <-ctx.Done():
// 			log.Printf("worker-%d shutting down", id)
// 			return
// 		}

// 	}

// }
func worker(id int, jobQueue <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobQueue {
		log.Printf("worker-%d processing job-%d", id, job.ID)
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("worker-%d stopped", id)
}
