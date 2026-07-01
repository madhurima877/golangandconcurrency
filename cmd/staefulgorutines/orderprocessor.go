package main

import "fmt"

type OrderProcessor struct {
	totalOrders int
	orders      chan string
}

func (o *OrderProcessor) START() {
	go func() {
		for order := range o.orders {
			o.totalOrders++
			fmt.Printf("Processing %s (Total: %d)\n", order, o.totalOrders)

		}
	}()
}
func (o *OrderProcessor) Sumbit(value string) {
	o.orders <- value

}
