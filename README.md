# Go FlexPool

[![GoDoc](https://godoc.org/github.com/pablolagosm/go-flexpool?status.svg)](https://godoc.org/github.com/pablolagosm/go-flexpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/pablolagosm/go-flexpool)](https://goreportcard.com/report/github.com/pablolagosm/go-flexpool)

Go FlexPool is a flexible and efficient goroutine pool implementation in Go, designed to manage and execute tasks with priority-based scheduling. It provides features such as dynamic resizing, context-based task cancellation, and error handling.

## Features

- Priority-based task scheduling
- Dynamic pool resizing
- Context-based task cancellation
- Graceful shutdown
- Error reporting through channels
- Configurable maximum tasks and workers

## Installation

To install the package, use the following command:

```sh
go get github.com/pablolagosm/go-flexpool
```

## Usage
###Creating a Pool
To create a new pool, specify the maximum number of workers and the maximum number of tasks:

```go
package main

import (
    "context"
    "fmt"
    "github.com/pablolagosm/go-flexpool"
    "time"
)

func main() {
    pool := flexpool.New(5, 10) // 5 workers, 10 maximum tasks

    ctx := context.Background()

    // Submitting tasks
    for i := 0; i < 15; i++ {
        index := i
        err := pool.submit(ctx, func(ctx context.Context) error {
            fmt.Printf("Executing task %d\n", index)
            time.Sleep(1 * time.Second)
            return nil
        }, flexpool.MediumPriority)
        
        if err != nil {
            fmt.Printf("Failed to submit task %d: %v\n", index, err)
        }
    }

    // Shutdown the pool
    shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    err := pool.Shutdown(shutdownCtx)
    if err != nil {
        fmt.Printf("Failed to shutdown pool: %v\n", err)
    }
}
```

### Submitting Tasks
Tasks can be submitted to the pool with a specific priority:

```go
err := pool.Submit(ctx, func(ctx context.Context) error {
// Task logic here
return nil
}, flexpool.HighPriority)
```

### Resizing the Pool
You can resize the pool dynamically:

```go
resizeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

err := pool.Resize(resizeCtx, 10) // Resize to 10 workers
if err != nil {
fmt.Printf("Failed to resize pool: %v\n", err)
}
```

### Shutdown the Pool
To gracefully shutdown the pool, ensuring all tasks are completed or cancelled:

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := pool.Shutdown(shutdownCtx)
if err != nil {
fmt.Printf("Failed to shutdown pool: %v\n", err)
}
```

### Error Handling
Errors encountered during task execution can be retrieved from the pool's error channel:

```go
select {
case err := <-pool.errorChan:
fmt.Printf("Task error: %v\n", err)
default:
// No error
}
```

### Contributing
Contributions are welcome! Please submit issues or pull requests for any enhancements or bug fixes.

License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements

This project was inspired by various goroutine pool implementations and aims to provide a flexible and robust solution for managing concurrent tasks in Go.

---

Feel free to adjust any sections or add more details as necessary for your project.