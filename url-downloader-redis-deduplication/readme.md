# URL Downloader with Redis Deduplication

A concurrent URL downloader built with Go that uses Redis to prevent duplicate downloads.

## Features

- Reads URLs from a file
- Concurrent worker pool with 10 workers
- Buffered job queue using Go channels
- HTTP URL downloading
- Per-request timeout using `context.WithTimeout`
- HTTP status validation
- Redis-based duplicate prevention
- Atomic URL claiming using Redis `SETNX`
- Exponential retry with backoff
- Random jitter between retries
- Graceful shutdown using context cancellation
- Worker synchronization using `sync.WaitGroup`

## How It Works

```text
urls.txt
   |
   v
Producer
   |
   v
Jobs Channel
   |
   v
Worker Pool
   |
   v
Redis SETNX
   |
   +--> Already claimed -> Skip
   |
   v
Download URL
   |
   +--> Success -> Mark as done in Redis
   |
   +--> Failure -> Retry with backoff and jitter
                      |
                      v
                 Retries exhausted
                      |
                      v
                Delete Redis key