package main

import (
	"golangandconcurrency/cmd/kafkagolang/kafkareader"
	"golangandconcurrency/cmd/kafkagolang/kafkawriter"
	"strconv"
	"sync"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafkawriter.NewWriter()

	defer writer.Close()
	var data []string
	for i := 0; i <= 100; i++ {
		data = append(data, strconv.Itoa(i))
	}

	kafkawriter.Writer(writer, data...)
	reader := kafkareader.NewReader()
	defer reader.Close()
	jobs := make(chan kafka.Message, 100)
	offsetchan := make(chan kafka.Message)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go kafkareader.Worker(reader, i, jobs, &wg, offsetchan)

	}
	go kafkareader.CommitManager(offsetchan, reader)

	kafkareader.KeepReading(reader, jobs)
	wg.Wait()

}
