package main

import (
	"log"
	"sync"
)

func semaphore() {
	nums := []int{1, 2, 3, 3, 4, 5, 6, 7, 8, 9, 10, 22, 12, 34, 56, 77, 88, 99}
	sem := make(chan struct{}, 3)

	var wg sync.WaitGroup
	for i := 0; i < len(nums); i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			processJob(i)
		}(i, &wg)
	}
	wg.Wait()

}
func processJob(i int) {
	log.Println(i)
}
