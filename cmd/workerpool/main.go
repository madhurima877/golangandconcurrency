package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}
	go func() {
		for i := 1; i <= 20; i++ {
			jobs <- i

		}
		

	}()
	wg.Wait()

}
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println("processing job", id, " job ", job)
		time.Sleep(1 * time.Second)
	}
}
