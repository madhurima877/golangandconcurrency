package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {

	jobs := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	shut := make(chan os.Signal, 1)
	signal.Notify(shut, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer(i, jobs, &wg)
	}

	go producer(ctx, jobs, "data.txt")

	<-shut
	fmt.Println("\nShutdown signal received")

	cancel()

	wg.Wait()

	fmt.Println("Application exited gracefully")
}

func consumer(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range jobs {
		fmt.Printf("Worker %d processing: %s\n", id, line)
		time.Sleep(2 * time.Second)
	}
	fmt.Printf("Worker %d exited\n", id)
}

func producer(ctx context.Context, jobs chan<- string, file string) {
	defer close(jobs)
	data, err := os.ReadFile("/var/www/html/gocode/src/golangandconcurrency/cmd/gracefulproducerconsumer/data.txt")
	if err != nil {
		fmt.Println("error reading file:", err)
		return
	}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		select {
		case <-ctx.Done():
			fmt.Println("Producer stopped")
			return
		case jobs <- line:
		}
	}

}
