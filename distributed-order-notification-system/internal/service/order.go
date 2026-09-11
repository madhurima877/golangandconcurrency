package service

import (
	"fmt"
	"sync"
	"time"
)

var (
	processedOrders int
	mu              sync.Mutex
)

func ProcessOrder(orderID int) {

	mu.Lock()
	processedOrders++
	mu.Unlock()

	fmt.Println("Processing order:", orderID)

	time.Sleep(2 * time.Second)

	fmt.Println("Order processed:", orderID)
}
func RunWorkers() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Println("Worker started:", id)
			time.Sleep(1 * time.Second)
			fmt.Println("Worker finished:", id)

		}(i)
	}
	wg.Wait()
	fmt.Println("All workers completed")
}
