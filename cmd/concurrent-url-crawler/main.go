package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	ID  int
	URL string
}

type Result struct {
	ID           int
	InputURL     string
	StatusCode   int
	ResponseTime time.Duration
	ContentLen   int64
	FinalURL     string
	Error        string
}

const (
	workerCount = 50
	maxRetries  = 3
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var workers sync.WaitGroup
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i := 0; i < workerCount; i++ {
		workers.Add(1)
		go Worker(ctx, i, &workers, jobs, results, client)
	}
	go func() {
		defer close(jobs)
		file, err := os.Open("urls.txt")
		if err != nil {
			cancel()
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		id := 0
		for scanner.Scan() {
			url := scanner.Text()
			select {
			case jobs <- Job{
				ID:  id,
				URL: url,
			}:
				id++
			case <-ctx.Done():
				return
			}
		}
		if err := scanner.Err(); err != nil {
			cancel()
		}
	}()

	go func() {
		workers.Wait()
		close(results)
	}()
	file, err := os.Create("results.csv")
	if err != nil {
		cancel()
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"id",
		"input_url",
		"status_code",
		"response_time",
		"content_length",
		"final_url",
		"error"})
	for result := range results {
		writer.Write([]string{
			strconv.Itoa(result.ID),
			result.InputURL,
			strconv.Itoa(result.StatusCode),
			result.ResponseTime.String(),
			strconv.FormatInt(result.ContentLen, 10),
			result.FinalURL,
			result.Error,
		})
	}
}
func Worker(ctx context.Context, workerID int, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result, client *http.Client) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			result := crawl(ctx, client, job)
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}

}
func crawl(ctx context.Context, client *http.Client, job Job) Result {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
		if err != nil {
			return Result{
				ID:       job.ID,
				InputURL: job.URL,
				Error:    err.Error(),
			}
		}
		resp, err := client.Do(req)
		elapsed := time.Since(start)
		if err == nil {
			defer resp.Body.Close()
			return Result{
				ID:           job.ID,
				InputURL:     job.URL,
				StatusCode:   resp.StatusCode,
				ResponseTime: elapsed,
				ContentLen:   resp.ContentLength,
				FinalURL:     resp.Request.URL.String(),
			}
		}
		lastErr = err
		if ctx.Err() != nil {
			return Result{
				ID:       job.ID,
				InputURL: job.URL,
				Error:    ctx.Err().Error(),
			}
		}
		delay := time.Duration(1<<uint(attempt-1)) * time.Second
		timer := time.NewTimer(delay)

		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return Result{
				ID:       job.ID,
				InputURL: job.URL,
				Error:    ctx.Err().Error(),
			}

		}
	}
	return Result{
		ID:       job.ID,
		InputURL: job.URL,
		Error:    lastErr.Error(),
	}
}
