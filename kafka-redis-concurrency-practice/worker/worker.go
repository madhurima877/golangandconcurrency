package worker

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"time"
)

type Worker struct {
	redis *redis.Client
	id    string
}

func New(redisClient *redis.Client, id string) *Worker {
	return &Worker{
		redis: redisClient,
		id:    id,
	}
}
func (w *Worker) Start(ctx context.Context, stream string, group string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
		streams, err := w.redis.ReadGroup(ctx, stream, group, w.id)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		for _, streamData := range streams {
			for _, msg := range streamData.Messages {
				fmt.Printf("Worker %s processing: %v\n", w.id, msg.Values)
				time.Sleep(2 * time.Second)

				err := w.redis.Ack(ctx, stream, group, msg.ID)
				if err != nil {
					fmt.Println("ack error:", err)
					continue
				}
				fmt.Printf(
					"worker %s completed: %s\n",
					w.id,
					msg.ID,
				)
			}
		}
	}
}
