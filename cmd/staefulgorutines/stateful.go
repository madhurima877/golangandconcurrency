package main

import (
	"log"
	"sync"
	"sync/atomic"
)

// type StatefulWorker struct {
// 	count int
// 	ch    chan int
// }

// func (w *StatefulWorker) increment() {
// 	go func() {
// 		for {
// 			// select {
// 			// case
// 			value := <-w.ch
// 			w.count += value
// 			log.Println("Current Value", w.count)
// 			// }
// 		}
// 	}()

// }
// func (w *StatefulWorker) send(value int) {
// 	w.ch <- value
// }
// func main() {
// 	stWorker := &StatefulWorker{ch: make(chan int)}
// 	stWorker.increment()

// 	for i := range 5 {
// 		stWorker.send(i)
// 		time.Sleep(500 * time.Millisecond)
// 	}
// 	processor := &OrderProcessor{
// 		orders: make(chan string),
// 	}
// 	processor.START()
// 	go func() {
// 		processor.orders <- "Laptop"
// 	}()
// 	go processor.Sumbit("LAPTOP2")
// 	go processor.Sumbit("Phone")
// 	go processor.Sumbit("Phone3")
// 	go processor.Sumbit("Mouse")
// 	time.Sleep(time.Second)
// }

type CountVal interface {
	increment()
	getValue() int64
}

type Counter struct {
	Num int64
}

func (C *Counter) increment() {
	atomic.AddInt64(&C.Num, 1)
}
func (c *Counter) getValue() int64 {
	return atomic.LoadInt64(&c.Num)
}

func main() {
	var c CountVal = &Counter{Num: 0}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				c.increment()
			}
		}()
	}

	wg.Wait()
	log.Println(c.getValue())
}
