package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			fmt.Println("Processing job...")
			time.Sleep(2 * time.Second)
		}
	}()

	<-stop
	fmt.Println("Received shutdown signal")
	fmt.Println("Closing database connections...")
	fmt.Println("Stopping workers...")
	fmt.Println("Application exited gracefully")
}
