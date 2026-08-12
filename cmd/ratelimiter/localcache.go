package main

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Cache struct {
	local map[string]string
	mu    sync.RWMutex
	group singleflight.Group
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.local[key]
	if !ok {
		return "", false
	}
	return val, true
}
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.local[key] = value
}
func (c *Cache) GetOrFetch(key string) (string, error) {
	val, ok := c.Get(key)
	if ok {
		return val, nil
	}
	result, err, _ := c.group.Do(key, func() (interface{}, error) {
		value, err := fetchFromRedis(key)
		if err != nil {
			return nil, err
		}

		c.Set(key, value)
		return value, nil
	})
	if err != nil {
		return "", err
	}

	return result.(string), nil

}
func fetchFromRedis(key string) (string, error) {
	time.Sleep(1 * time.Second)
	return "value-from-redis", nil
}

func main() {

}
