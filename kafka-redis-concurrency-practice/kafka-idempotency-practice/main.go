package main

import (
	"context"
	"errors"
	"fmt"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"golangandconcurrency/kafka-redis-concurrency-practice/postgres"
	"log"
	"time"

	"github.com/lib/pq"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	db, err := postgres.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println("PostgreSQL connected")

	ctx := context.Background()
	writer := kafka.NewWriter("idempotency-practice")
	defer writer.Close()

	eventID := fmt.Sprintf("event-%d", time.Now().UnixNano())
	for i := 0; i < 2; i++ {
		err := writer.WriteMessageWithKey(ctx, []byte(eventID), []byte("order-created"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("produced event: %s", eventID)
	}

	reader := kafka.NewReader("idempotency-practice",
		"idempotency-group")
	defer reader.Close()
	for i := 0; i < 2; i++ {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("consumed event: key=%s value=%s",
			msg.Key, msg.Value)
		_, err = db.ExecContext(
			ctx,
			"INSERT INTO processed_events (event_id) VALUES ($1)",
			string(msg.Key),
		)
		if err != nil {
			var pqErr *pq.Error

			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				log.Printf("duplicate event ignored: %s", msg.Key)

				if err := reader.CommitMessage(ctx, msg); err != nil {
					log.Fatal(err)
				}

				continue
			}

			log.Printf(
				"failed to process event %s: %v",
				msg.Key,
				err,
			)

			continue
		}
		log.Printf("event processed: %s", msg.Key)

		// log.Fatal("SIMULATED CRASH BEFORE KAFKA COMMIT")

		if err := reader.CommitMessage(ctx, msg); err != nil {
			log.Fatal(err)
		}
	}
}
