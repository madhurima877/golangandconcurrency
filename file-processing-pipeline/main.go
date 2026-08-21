package main

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fileJobs := make(chan string, 100)
	lines := make(chan string, 100)
	results := make(chan string, 100)
	var fileWorkers sync.WaitGroup
	for i := 1; i <= 10; i++ {
		fileWorkers.Add(1)
		go FileWorker(ctx, &fileWorkers, fileJobs, lines)
	}
	go func() {
		entries, err := os.ReadDir("logs")
		if err != nil {
			cancel()
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			filepath := filepath.Join("logs", entry.Name())
			select {
			case <-ctx.Done():
				return
			case fileJobs <- filepath:
			}

		}
		close(fileJobs)
	}()
	go func() {
		fileWorkers.Wait()
		close(lines)
	}()
	var processWorkers sync.WaitGroup
	for i := 1; i <= 20; i++ {
		processWorkers.Add(1)
		go ProcessWorker(ctx, &processWorkers, lines, results)
	}
	go func() {
		processWorkers.Wait()
		close(results)
	}()
	file, err := os.Create("output.txt")
	if err != nil {
		return
	}
	defer file.Close()
	for res := range results {
		file.Write([]byte(res + "\n"))
	}
}

func FileWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, lines chan<- string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			file, err := os.Open(job)
			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(file)
			defer file.Close()
			for scanner.Scan() {
				line := scanner.Text()
				select {
				case <-ctx.Done():
					return
				case lines <- line:
				}

			}
			if err := scanner.Err(); err != nil {
				continue
			}

		}
	}
}
func ProcessWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			processed := strings.ToUpper(job)
			select {
			case <-ctx.Done():
				return
			case results <- processed:
			}
		}

	}
}
