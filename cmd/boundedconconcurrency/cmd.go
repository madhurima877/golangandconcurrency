package main

import (
	"errors"
	"log"
	"sync"
)

func main() {
	sem := make(chan struct{}, 5)
	nums := []int{1, 2, 3, 3, 4, 5, 6, 7, 8, 9, 10, 22, 12, 34, 56, 77, 88, 99}
	var wg sync.WaitGroup
	errCh := make(chan error, len(nums))
	for i := 0; i < len(nums); i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			if err := ProcessOrders(i); err != nil {
				errCh <- err
			}

		}(nums[i])

	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		log.Println(err)
	}
}
func ProcessOrders(orderIDs int) error {
	log.Println(orderIDs)
	if orderIDs == 3 {
		return errors.New("error detected")
	}
	return nil
}
