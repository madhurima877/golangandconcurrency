package service

import (
	"fmt"
	"sync"
	"time"
)

var (
	processedOrders int
	mu              sync.Mutex
	workerSemaphore = make(chan struct{}, 3)
)

func ProcessOrder(orderID int) {
	workerSemaphore <- struct{}{}
	defer func() {
		<-workerSemaphore
	}()

	mu.Lock()
	processedOrders++
	mu.Unlock()

	fmt.Println("Processing order:", orderID)

	time.Sleep(2 * time.Second)

	fmt.Println("Order processed:", orderID)
}
