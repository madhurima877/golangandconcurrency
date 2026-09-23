package main

import "sync"

func main() {
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i := 0; i <= 10; i++ {
		<-sem
		wg.Add(1)
		go func() {
			defer func() {
				sem <- struct{}{}
			}()
			defer wg.Done()
		}()
	}
	wg.Wait()
}
