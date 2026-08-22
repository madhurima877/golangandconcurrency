package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const maxRetries = 3

type RedisCleint struct {
	client *redis.Client
}

func NewRedisClient() *RedisCleint {
	client := redis.NewClient(&redis.Options{Addr: "localhost:16390"})
	return &RedisCleint{client: client}
}
func (r *RedisCleint) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisClient := NewRedisClient()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	jobs := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go redisClient.DownloadWorker(ctx, jobs, &wg)
	}

	go ProduceURLs(ctx, jobs)
	WaitForShutdown(cancel)
	wg.Wait()
	log.Println("all jobs cancelled")
}
func ProduceURLs(ctx context.Context, jobs chan<- string) {
	defer close(jobs)
	file, err := os.Open("urls.txt")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		select {
		case <-ctx.Done():
			return
		case jobs <- line:

		}

	}
	if err := scanner.Err(); err != nil {
		return
	}

}

func (r *RedisCleint) DownloadWorker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			retries := 0

			key := "job:" + job
			claimed, err := r.client.SetNX(ctx, key, "processing", 10*time.Second).Result()

			if err != nil {
				log.Printf("redis error for %s: %v", job, err)
				continue
			}

			if !claimed {
				log.Printf("skipping duplicate: %s", job)
				continue
			}
			log.Printf("claimed URL: %s", job)

			func() {
				success := false
				for retries < maxRetries {
					ctxx, cancel := context.WithTimeout(ctx, 5*time.Second)

					req, err := http.NewRequestWithContext(ctxx, http.MethodGet, job, nil)
					if err != nil {
						log.Printf("download failed for %s: %v", job, err)
						cancel()
						retries++
						if retries < maxRetries {
							RetryDelay(ctx, retries)
						}

						continue
					}
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Printf("download failed for %s: %v", job, err)
						cancel()
						retries++
						if retries < maxRetries {
							RetryDelay(ctx, retries)
						}
						continue
					}

					if resp.StatusCode < 200 || resp.StatusCode > 299 {
						log.Printf("download failed for %s: status %d", job, resp.StatusCode)

						cancel()
						retries++
						if retries < maxRetries {
							RetryDelay(ctx, retries)
						}
						resp.Body.Close()
						continue
					}

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Printf("download failed for %s: %v", job, err)
						cancel()
						retries++
						if retries < maxRetries {
							RetryDelay(ctx, retries)
						}
						resp.Body.Close()
						continue
					}
					log.Printf("downloaded %s, bytes: %d", job, len(body))
					resp.Body.Close()
					success = true
					cancel()
					break

				}
				if success {
					if err := r.client.Set(ctx, key, "done", 10*time.Second).Err(); err != nil {
						log.Printf("failed to mark %s as done: %v", job, err)
						return
					}

					log.Printf("marked as downloaded: %s", job)
				} else {
					log.Printf("retries exhausted for %s", job)
					r.client.Del(ctx, key)
				}

			}()

		}
	}
}
func WaitForShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	cancel()
}
func RetryDelay(ctx context.Context, retryCount int) {
	randomJitter := rand.Intn(9-1) + 1
	delay := time.Duration(1<<(retryCount-1)) * time.Second
	delay += time.Duration(randomJitter) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
	}
}
