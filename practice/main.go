package main

import (
	"context"
	"log"
	"sync"
)

func main() {
	arr := []int{1, 2}
	appends(&arr, 3)
	log.Println(arr)
}
func appends(arr *[]int, i int) {
	*arr = append(*arr, i)
}

func processWithWorkers(ctx context.Context, tasks []func() error, workerCount int) []error {
	var wg sync.WaitGroup
	res := make(chan error)
	result := []error{}
	jobs := make(chan func() error)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(&wg, jobs, ctx, res)
	}

	go func() {
		defer close(jobs)
		for i := 0; i < len(tasks); i++ {
			jobs <- tasks[i]

		}
	}()
	go func() {

		wg.Wait()
		close(res)
	}()

	for {
		select {
		case <-ctx.Done():
			return result
		case d, ok := <-res:
			if !ok {
				return result
			}
			result = append(result, d)
		}

	}
	//return result
}
func worker(wg *sync.WaitGroup, jobs <-chan func() error, ctx context.Context, res chan<- error) {
	defer wg.Done()
	for {

		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			err := job()

			res <- err
		}

	}
}
