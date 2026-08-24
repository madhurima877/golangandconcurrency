# File Redis Pipeline

A simple asynchronous file processing pipeline built with Go, Gin, Redis Streams, and Docker.

## Architecture

Client
  ↓
Gin API
  ↓
Save File
  ↓
Redis Stream (XADD)
  ↓
Worker Consumer Group (XREADGROUP)
  ↓
Process File
  ↓
XACK