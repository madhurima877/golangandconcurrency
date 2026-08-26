package fileworker

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"time"
)

type Worker struct {
	redisClient *redis.Client
	id          string
}

func New(redisClient *redis.Client, id string) *Worker {
	return &Worker{
		redisClient: redisClient,
		id:          id,
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		job, err := w.redisClient.PopJob(ctx, "file-jobs")
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		fmt.Printf("worker %s processing file: %s\n", w.id, job)

		time.Sleep(2 * time.Second)

		fmt.Printf("worker %s finished file: %s\n", w.id, job)
	}
}
