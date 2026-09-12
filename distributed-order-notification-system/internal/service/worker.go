package service

import (
	"log"
	"sync"
	"time"
)

type Job struct {
	OrderID int
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		log.Println("Worker", id, "processing order", job.OrderID)
		time.Sleep(2 * time.Second)
		log.Println("Worker", id, "finished order", job.OrderID)
	}
}

func StartWorkerPool(workerCount int, jobs <-chan Job, wg *sync.WaitGroup) {

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, wg)
	}

}

func SubmitJob(jobs chan<- Job, orderID int) {
	jobs <- Job{OrderID: orderID}
}

func RunWorkerPool() {
	var wg sync.WaitGroup
	jobs := make(chan Job)
    StartWorkerPool(3, jobs, &wg)
	for i := 1; i <= 10; i++ {
		SubmitJob(jobs, i)
	}
	close(jobs)
	wg.Wait()

}
