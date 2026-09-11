package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	cilent := redis.NewClient(&redis.Options{
		Addr: "locahost:63480",
	})
	return &RedisClient{client: cilent}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:63480",
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil

}
func (r *RedisClient) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
