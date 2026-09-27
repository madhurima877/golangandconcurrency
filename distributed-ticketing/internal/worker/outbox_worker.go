package worker

import (
	"database/sql"
	"fmt"
)

type OutboxWorker struct {
	db *sql.DB
}

func NewOutboxWorker(db *sql.DB) *OutboxWorker {
	return &OutboxWorker{db: db}
}

func (w *OutboxWorker) Process() error {
	rows, err := w.db.Query(`SELECT id, event_type, aggregate_id, payload
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY id
		LIMIT 100`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var eventType string
		var aggregateID int64
		var payload []byte
		if err := rows.Scan(&id, &eventType, &aggregateID, &payload); err != nil {
			return err
		}
		fmt.Printf("Outbox event: id=%d type=%s aggregateID=%d payload=%s\n",
			id, eventType, aggregateID, payload)
	}
	return rows.Err()
}
