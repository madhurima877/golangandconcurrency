package kafka

import kafkago "github.com/segmentio/kafka-go"

type Reader struct {
	reader *kafkago.Reader
}

type Writer struct {
	writer *kafkago.Writer
}

func NewReader(topic, groupID string) *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: groupID,
		Topic:   topic,
	})

	return &Reader{
		reader: reader,
	}
}

func NewWriter(topic string) *Writer {
	writer := &kafkago.Writer{
		Addr:     kafkago.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}

	return &Writer{
		writer: writer,
	}
}
