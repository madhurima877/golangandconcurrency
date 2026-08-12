package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rbd := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	for {
		key := "user:1015"
		value, err := rbd.Get(ctx, key).Result()
		if err == redis.Nil {
			log.Println("CACHE MISS")

			value = "MAdhurima"

			err = rbd.Set(ctx, key, value, 30*time.Second).Err()
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(nil)
		} else {
			fmt.Println("CACHE HIT")
		}
		fmt.Println("Value:", value)
	}

}
