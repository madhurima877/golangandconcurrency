package main

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	workerCount int
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{workerCount: workerCount}
}

func (wp *WorkerPool) Start(jobs <-chan string, done *sync.WaitGroup) {
	defer done.Done()
	var wg sync.WaitGroup
	for i := 1; i <= wp.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("worker-%d started %s\n", workerID, job)

				time.Sleep(2 * time.Second)

				fmt.Printf("worker-%d finished %s\n", workerID, job)
			}
		}(i)
	}
	wg.Wait()
}
