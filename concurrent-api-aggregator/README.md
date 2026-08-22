# Concurrent API Aggregator

A Go application that concurrently calls multiple external APIs and aggregates their responses into a single JSON response.

## Features

- Concurrent API requests using goroutines
- Fan-out and fan-in pattern
- WaitGroup synchronization
- Channel-based result collection
- Partial failure handling
- Overall request timeout
- Per-API timeout
- Context cancellation
- HTTP status code validation
- Concurrency limiting using a buffered channel semaphore

## How It Works

The `/aggregate` endpoint starts multiple API requests concurrently.

A semaphore limits the number of simultaneous requests.

Each API request has its own timeout, while the entire aggregation request also has an overall timeout.

Successful and failed responses are collected through a channel and returned as a single JSON response.

## Run

```bash
go run .