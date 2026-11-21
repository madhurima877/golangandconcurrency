package main

import (
	"context"
	"errors"
	"fmt"
)

type Result struct {
	URL  string
	Body []byte
	Err  error
}

var ErrTooManyFailures = errors.New("too many failures")

func FetchAll(ctx context.Context, urls []string, maxConcurrent int, fetch func(ctx context.Context, url string) ([]byte, error)) ([]Result, error) {
	if maxConcurrent <= 0 {
		return nil, fmt.Errorf("fetchall:maxConconcurrent must be >0,got %d", maxConcurrent)
	}
	if len(urls) == 0 {
		return []Result{}, nil
	}
	result := make([]Result, len(urls))
}
