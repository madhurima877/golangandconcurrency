package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis *redis.Client
}

func New() *Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Client{
		Redis: client,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Redis.Ping(ctx).Err()
}
