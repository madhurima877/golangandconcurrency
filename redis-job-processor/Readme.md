# Redis-backed Concurrent Job Processor

## Overview

A concurrent job processing system built with Go and Redis.

Jobs are read from a file, pushed into a Redis queue, and processed by multiple concurrent consumers.

## Architecture

```text
jobs.txt
   |
   v
Producer
   |
   v
Redis Queue ("jobs")
   |
   v
20 Consumers
   |
   v
Semaphore
(Max 5 concurrent jobs)
   |
   +--> Success
   |
   +--> Failure
          |
          v
       Retry
    1s -> 2s -> 4s
          |
          v
        Redis Queue
          |
          v
      Max Retries
          |
          v
     Redis DLQ
    ("jobs:dlq")