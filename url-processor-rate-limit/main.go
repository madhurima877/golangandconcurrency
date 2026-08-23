package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	var queued int64
	var completed int64
	var failed int64
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		log.Println("shutdown signal received")
	}()
	limiter := time.NewTicker(500 * time.Millisecond)
	defer limiter.Stop()
	urls := []string{
		"https://example.com",
		"https://www.google.com",
		"https://httpbin.org/get",
		"https://example.org",
		"https://httpbin.org/delay/5",
	}
	jobs := make(chan string)

	log.Println("URL Processor + Global Rate Limit")
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go Worker(i, jobs, &wg, ctx, limiter, &completed, &failed)
	}

	go Producer(urls, jobs, ctx, &queued)
	wg.Wait()
	log.Println("----- Metrics -----")
	log.Printf("queued: %d", atomic.LoadInt64(&queued))
	log.Printf("completed: %d", atomic.LoadInt64(&completed))
	log.Printf("failed: %d", atomic.LoadInt64(&failed))
}
func Producer(urls []string, jobs chan<- string, ctx context.Context, queued *int64) {
	defer close(jobs)
	for _, url := range urls {
		select {
		case <-ctx.Done():
			return
		case jobs <- url:
			atomic.AddInt64(queued, 1)
			log.Printf("queued URL: %s", url)
		}
	}
}
func Worker(i int, jobs <-chan string, wg *sync.WaitGroup, ctx context.Context, limiter *time.Ticker, completed, failed *int64) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			log.Printf("worker %d received URL: %s", i, job)

			select {
			case <-ctx.Done():
				return
			case <-limiter.C:
				log.Printf("worker %d allowed by rate limiter", i)
			}
			reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

			req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, job, nil)
			if err != nil {
				atomic.AddInt64(failed, 1)
				cancel()
				log.Printf("worker %d failed reading response: %v", i, err)
				continue
			}

			func() {
				defer cancel()
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					atomic.AddInt64(failed, 1)
					log.Printf("worker %d request failed: %v", i, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode < 200 || resp.StatusCode > 299 {
					log.Printf(
						"worker %d received status %d for %s",
						i,
						resp.StatusCode,
						job,
					)
					atomic.AddInt64(failed, 1)
					return
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("worker %d failed reading response: %v", i, err)
					atomic.AddInt64(failed, 1)
					return
				}
				atomic.AddInt64(completed, 1)
				log.Printf(
					"worker %d completed URL: %s, status: %d, bytes: %d",
					i,
					job,
					resp.StatusCode,
					len(body),
				)
			}()

		}
	}
}
