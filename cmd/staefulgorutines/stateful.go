package main

import (
	"log"
	"time"
)

type StatefulWorker struct {
	count int
	ch    chan int
}

func (w *StatefulWorker) increment() {
	go func() {
		for {
			// select {
			// case
			value := <-w.ch
			w.count += value
			log.Println("Current Value", w.count)
			// }
		}
	}()

}
func (w *StatefulWorker) send(value int) {
	w.ch <- value
}
func main() {
	stWorker := &StatefulWorker{ch: make(chan int)}
	stWorker.increment()

	for i := range 5 {
		stWorker.send(i)
		time.Sleep(500 * time.Millisecond)
	}
	processor := &OrderProcessor{
		orders: make(chan string),
	}
	processor.START()
	go func() {
		processor.orders <- "Laptop"
	}()
	go processor.Sumbit("LAPTOP2")
	go processor.Sumbit("Phone")
	go processor.Sumbit("Phone3")
	go processor.Sumbit("Mouse")
	time.Sleep(time.Second)
}
