package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	jobs := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go producer(ctx, jobs, &wg, i)
	}
	go consumer(jobs)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutdown signal received")
	cancel()
	wg.Wait()
	close(jobs)
	time.Sleep(3 * time.Second)

	fmt.Println("Application exited gracefully")

}
func producer(ctx context.Context, jobs chan<- string, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Producer %d stopped\n", id)
			return
		case jobs <- fmt.Sprintf("P%d-Job%d", id, i):
			time.Sleep(time.Second)
		}

	}
}

func consumer(jobs <-chan string) {
	for job := range jobs {
		fmt.Println("Processing:", job)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Consumer exited")
}
