package main

import (
	"context"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"log"
)

func Consume() {
	ctx := context.Background()
	reader := kafka.NewReader(
		"partition-practice",
		"partition-practice-group-v2",
	)
	defer reader.Close()

	for i := 0; i < 3; i++ {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"key=%s value=%s partition=%d offset=%d\n",
			string(msg.Key),
			string(msg.Value),
			msg.Partition,
			msg.Offset,
		)
		if err := reader.CommitMessage(ctx, msg); err != nil {
			log.Fatal(err)
		}
	}

}
