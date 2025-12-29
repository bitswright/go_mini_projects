# Project 3: Concurrent URL Status Checker

This project is a command-line tool written in Go that checks the availability and response time of multiple URLs concurrently.

It demonstrates how Go’s concurrency primitives (goroutines, channels, and worker pools) can be used to efficiently perform I/O-bound tasks such as HTTP requests.

## Requirements
- Check multiple URLs concurrently
- Display HTTP status code and response time
- Handle request failures gracefully
- Limit number of concurrent requests
- Configurable timeout
- Optional retry mechanism
- Example Usage
    ```bash
    go run main.go -file urls.txt -workers 5 -timeout 3

    https://google.com      200   120ms
    https://github.com      200   95ms
    https://bad-url.test    ERROR timeout
    ```
- Command-line Flags
    | Flag        | Description                            | Default  |
    | ----------- | -------------------------------------- | -------- |
    | -file       | File containing URLs (one per line)    | urls.txt |
    | -workers    | Number of concurrent workers           | 5        |
    | -timeout    | Request timeout (seconds)              | 3        |
    | -retries    | Number of retries per URL              | 0        |
