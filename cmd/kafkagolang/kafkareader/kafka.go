package kafkareader

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

func NewReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "consumer-group",
	})
}

func KeepReading(reader *kafka.Reader, jobs chan<- kafka.Message) {
	defer close(jobs)
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		// var msgstring string
		// if err := json.Unmarshal(msg.Value, &msgstring); err != nil {
		// 	log.Println(err)
		// 	continue
		// }
		// msgstring = string(msg.Value)
		jobs <- msg

	}

}

var smp []int

func Worker(reader *kafka.Reader, workerid int, jobs <-chan kafka.Message, wg *sync.WaitGroup, offsetchan chan<- kafka.Message) {
	defer wg.Done()
	for data := range jobs {
		log.Printf("Processing job%s from worker id %d", string(data.Value), workerid)
		offsetchan <- data

	}

}

var processed = map[int64]kafka.Message{}

var lastCommitted = int64(0)

func CommitManager(offsetChan <-chan kafka.Message, reader *kafka.Reader) {
	for msg := range offsetChan {
		processed[msg.Offset] = msg

		for {
			next := lastCommitted + 1

			msgToCommit, ok := processed[next]
			if !ok {
				break
			}

			if err := reader.CommitMessages(context.Background(), msgToCommit); err != nil {
				log.Println(err)
				break
			}

			delete(processed, next)
			lastCommitted = next
		}
	}
}
