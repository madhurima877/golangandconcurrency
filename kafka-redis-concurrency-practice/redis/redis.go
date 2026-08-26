package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis *redis.Client
}

func New() *Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})

	return &Client{
		Redis: client,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Redis.Ping(ctx).Err()
}
func (c *Client) AddMessage(ctx context.Context, stream string, values map[string]interface{}) error {
	return c.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Err()
}

func (c *Client) CreateConsumerGroup(ctx context.Context, stream string, group string) error {
	return c.Redis.XGroupCreateMkStream(ctx, stream, group, "0").Err()
}

func (c *Client) ReadGroup(ctx context.Context, stream string, group string, consumer string) ([]redis.XStream, error) {
	return c.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    1,
		Block:    0,
	}).Result()
}
func (c *Client) Ack(ctx context.Context, stream string, group string, messageID string) error {
	return c.Redis.XAck(ctx, stream, group, messageID).Err()
}
func (c *Client) PushJob(ctx context.Context, queue string, value string) error {
	return c.Redis.LPush(ctx, queue, value).Err()

}
func (c *Client) PopJob(ctx context.Context, queue string) (string, error) {
	result, err := c.Redis.BRPop(ctx, time.Second, queue).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}
