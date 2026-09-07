package main

import (
	"context"
	"golangandconcurrency/kafka-redis-concurrency-practice/redis"
	"log"
	"sync"
)

func main() {
	ctx := context.Background()
	redisClient := redis.New()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	redisClient.Redis.Set(ctx, "counter", 0, 0)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := redisClient.IncrementWithWatch(ctx)
			if err != nil {
				log.Println("increment failed:", err)
			}
		}()
	}
	wg.Wait()
	value, _ := redisClient.Redis.Get(ctx, "counter").Int()
	log.Println("final counter:", value)

}

// func main() {
// 	var counter int
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			mu.Lock()
// 			counter++
// 			mu.Unlock()
// 		}()
// 	}
// 	wg.Wait()
// 	log.Println("Go counter:", counter)
// }
