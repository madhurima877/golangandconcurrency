# File Processing Pipeline

A concurrent file processing pipeline built in Go.

The program reads multiple log files from a directory, processes their contents concurrently, and writes the processed results into a single output file.

## Architecture

```text
logs/
  ↓
fileJobs channel
  ↓
10 File Reader Workers
  ↓
lines channel
  ↓
20 Processing Workers
  ↓
results channel
  ↓
1 Writer
  ↓
output.txt
```

## Features

* Reads multiple files from the `logs/` directory
* Processes files concurrently using 10 file reader workers
* Reads files line by line without loading entire files into memory
* Processes lines concurrently using 20 processing workers
* Uses bounded channels for backpressure
* Uses `context` for cancellation
* Continues processing other files if one file fails
* Uses WaitGroups for worker lifecycle management
* Safely closes channels between pipeline stages
* Uses a single writer to write all processed results to `output.txt`

## Project Structure

```text
file-processing-pipeline/
├── main.go
├── go.mod
├── logs/
│   ├── file1.log
│   ├── file2.log
│   └── file3.log
└── output.txt
```

## Pipeline Flow

### 1. Directory Scanner

The program reads the `logs/` directory and sends each file path into the `fileJobs` channel.

### 2. File Reader Workers

10 workers receive file paths from `fileJobs`.

Each worker:

* Opens a file
* Reads it line by line using `bufio.Scanner`
* Sends each line into the `lines` channel
* Continues processing if another file fails

### 3. Processing Workers

20 workers receive lines from the `lines` channel.

Each line is processed by converting it to uppercase:

```text
hello world
```

becomes:

```text
HELLO WORLD
```

The processed result is then sent to the `results` channel.

### 4. Writer

A single writer receives processed results and writes them to:

```text
output.txt
```

## Running the Project

Add some log files:

```bash
echo "hello world" > logs/file1.log
echo "go concurrency" > logs/file2.log
echo "worker pipeline" > logs/file3.log
```

Run:

```bash
go run .
```

Check the output:

```bash
cat output.txt
```

Example output:

```text
HELLO WORLD
GO CONCURRENCY
WORKER PIPELINE
```

The order may vary because files and lines are processed concurrently.

## Concurrency Model

The application uses a staged producer-consumer pipeline:

```text
Directory Scanner
      ↓
   fileJobs
      ↓
10 File Workers
      ↓
     lines
      ↓
20 Process Workers
      ↓
    results
      ↓
 Single Writer
```

Each stage communicates through bounded channels, preventing unlimited memory growth and naturally applying backpressure when downstream stages are slower.

## Graceful Completion

* The directory scanner closes `fileJobs` when no more files remain.
* File workers exit when `fileJobs` is closed.
* After all file workers finish, `lines` is closed.
* Processing workers exit when `lines` is closed.
* After all processing workers finish, `results` is closed.
* The writer finishes after consuming all results.

This ensures that every worker has a clear lifecycle and the pipeline shuts down without goroutine leaks or deadlocks.
