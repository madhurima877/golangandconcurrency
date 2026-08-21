package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redisclient struct {
	cleint *redis.Client
}

func NewClient() *Redisclient {
	client := redis.NewClient(&redis.Options{Addr: "localhost:16389"})
	return &Redisclient{
		cleint: client,
	}
}
func (r *Redisclient) Ping(ctx context.Context) error {
	return r.cleint.Ping(ctx).Err()
}

type Job struct {
	ID         int
	Payload    string
	RetryCount int
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redisClient := NewClient()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	go ProduceJobs(ctx, redisClient)
	semaphore := make(chan struct{}, 5)
	for i := 1; i <= 20; i++ {

		go ConsumeJobs(ctx, redisClient, semaphore)
	}
	WaitForShutdown(cancel)
	log.Println("all shutting down")

}

func ProduceJobs(ctx context.Context, client *Redisclient) {
	file, err := os.Open("jobs.txt")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		job := Job{
			ID:         i,
			Payload:    line,
			RetryCount: 0,
		}
		select {
		case <-ctx.Done():
			return
		default:
			jobbyte, err := json.Marshal(job)
			if err != nil {
				log.Println(err)
				return
			}
			err = client.cleint.LPush(ctx, "jobs", jobbyte).Err()
			if err != nil {
				log.Println(err)
				return
			}
			i++
		}

	}
	if err := scanner.Err(); err != nil {
		return
	}
	log.Println("all jobs produced")
}

func ConsumeJobs(ctx context.Context, client *Redisclient, semaphore chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			val, err := client.cleint.BLPop(ctx, 0, "jobs").Result()
			if err != nil {
				continue
			}

			jobData := val[1]
			var jobstuct Job
			err = json.Unmarshal([]byte(jobData), &jobstuct)
			if err != nil {
				continue
			}
			select {
			case semaphore <- struct{}{}:
			case <-ctx.Done():
				return
			}
			log.Printf("consumer started processing job %d", jobstuct.ID)

			err = ProcessJob(jobstuct)
			if err != nil {
				log.Println("failed")
				jobstuct.RetryCount++
				if jobstuct.RetryCount <= 3 {
					RetryDelay(jobstuct.RetryCount)
					select {
					case <-ctx.Done():
						return
					default:
						jobbyte, _ := json.Marshal(jobstuct)
						client.cleint.LPush(ctx, "jobs", jobbyte)

					}
				} else {
					jobbyte, _ := json.Marshal(jobstuct)
					client.cleint.LPush(ctx, "jobs:dlq", jobbyte)
					log.Printf("job %d moved to DLQ", jobstuct.ID)
				}

			} else {
				log.Println("success")
			}
			time.Sleep(2 * time.Second)
			log.Printf("consumer finished job %d", jobstuct.ID)

			<-semaphore

		}
	}
}
func ProcessJob(job Job) error {
	if job.ID%3 == 0 {
		return errors.New("failed")
	}
	return nil
}
func RetryDelay(retryCount int) {
	time.Sleep(time.Duration(1<<(retryCount-1)) * time.Second)
}
func WaitForShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	cancel()

}
