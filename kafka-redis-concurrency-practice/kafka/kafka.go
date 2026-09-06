package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	reader *kafkago.Reader
}

type Writer struct {
	writer *kafkago.Writer
}

func NewReader(topic, groupID string) *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        []string{"localhost:9082"},
		GroupID:        groupID,
		Topic:          topic,
		CommitInterval: 0,
	})

	return &Reader{
		reader: reader,
	}
}

func NewWriter(topic string) *Writer {
	writer := kafkago.NewWriter(kafkago.WriterConfig{
		Brokers:  []string{"localhost:9082"},
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	})

	return &Writer{
		writer: writer,
	}
}
func (w *Writer) Close() error {
	return w.writer.Close()
}

func (r *Reader) CommitMessage(ctx context.Context, msg kafkago.Message) error {
	return r.reader.CommitMessages(ctx, msg)

}
func (w *Writer) WriteMessage(ctx context.Context, value []byte) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{Value: value})
}

func (w *Writer) WriteMessageWithKey(ctx context.Context, key []byte, value []byte) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{
		Key:   key,
		Value: value,
	})
}
func (r *Reader) Close() error {
	return r.reader.Close()
}

func (w *Writer) WriteMessageWithKeyRoundRobin(ctx context.Context, key []byte, value []byte) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{
		Key:   key,
		Value: value,
	})
}
func ReaderWithWorker(ctx context.Context, topic string, groupID, workerId string) {
	reader := NewReader(topic, groupID)
	defer reader.Close()
	log.Printf("worker %s started", workerId)
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("worker %s shutting down", workerId)
				return
			}
			log.Printf("worker %s error: %v", workerId, err)
			continue
		}
		log.Printf(
			"worker=%s partition=%d key=%s value=%s",
			workerId,
			msg.Partition,
			msg.Key,
			msg.Value,
		)
		if err := reader.CommitMessage(ctx, msg); err != nil {
			log.Printf("worker %s commit error: %v", workerId, err)
		}

	}

}
func NewKeyWriter(topic string) *Writer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9082"),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	return &Writer{writer: writer}
}

func (r *Reader) FetchMessage(ctx context.Context) (kafkago.Message, error) {
	return r.reader.FetchMessage(ctx)
}
