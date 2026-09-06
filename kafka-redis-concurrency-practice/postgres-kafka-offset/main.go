package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/kafka"
	"golangandconcurrency/kafka-redis-concurrency-practice/postgres"
	"log"

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

	reader := kafka.NewReader("postgres-offset-practice",
		"postgres-offset-group")
	defer reader.Close()
	log.Println("Kafka consumer started")

	ctx := context.Background()
	msg, err := reader.FetchMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"received: partition=%d offset=%d value=%s",
		msg.Partition,
		msg.Offset,
		msg.Value,
	)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec(
		"INSERT INTO orders (order_id, status) VALUES ($1, $2)",
		string(msg.Value),
		"processed",
	)
	if err != nil {
		tx.Rollback()
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("duplicate order %s, skipping", msg.Value)

			if err := reader.CommitMessage(ctx, msg); err != nil {
				log.Printf("failed to commit duplicate offset: %v", err)
				return
			}

			log.Printf("duplicate offset committed")
			return
		}
		return
	}
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
	log.Printf("DB transaction committed for %s", msg.Value)
	if err := reader.CommitMessage(ctx, msg); err != nil {
		log.Printf("Kafka offset commit failed: %v", err)
		return
	}
	log.Printf("Kafka offset committed for %s", msg.Value)
}
