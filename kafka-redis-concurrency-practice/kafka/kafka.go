package kafka

import (
	"context"

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

func (r *Reader) FetchMessage(ctx context.Context) (kafkago.Message, error) {
	return r.reader.FetchMessage(ctx)
}
func (r *Reader) CommitMessage(ctx context.Context, msg kafkago.Message) error {
	return r.reader.CommitMessages(ctx, msg)
	return nil
}
func (w *Writer) WriteMessage(ctx context.Context, value []byte) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{Value: value})
}
