package kafkawriter

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func NewWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})
}
func Writer(writer *kafka.Writer, data ...string) {
	var messages []kafka.Message
	for _, d := range data {
		messages = append(messages, kafka.Message{
			Value: []byte(d),
		})
	}

	err := writer.WriteMessages(context.Background(), messages...)
	if err != nil {
		log.Println(err)
	}

}
