package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16380",
	})

	return &RedisCache{cache: client}
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.cache.Ping(ctx).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.cache.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key, value string) error {
	return r.cache.Set(ctx, key, value, 5*time.Minute).Err()
}
