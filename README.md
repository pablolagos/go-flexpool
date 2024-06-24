# Go FlexPool

[![GoDoc](https://godoc.org/github.com/pablolagos/go-flexpool?status.svg)](https://godoc.org/github.com/pablolagos/go-flexpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/pablolagos/go-flexpool)](https://goreportcard.com/report/github.com/pablolagos/go-flexpool)

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
go get github.com/pablolagos/go-flexpool
```

## Usage
### Creating a Pool
To create a new pool, specify the maximum number of workers and the maximum number of tasks:

```go
package main

import (
	"context"
	"fmt"
	"github.com/pablolagos/go-flexpool"
	"time"
)

func main() {
	// Create a new pool with 5 workers and a maximum of 10 tasks.
	pool := flexpool.New(5, 10)

	// Create a context with a timeout for submitting tasks.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Submit tasks to the pool.
	for i := 0; i < 15; i++ {
		index := i
		err := pool.Submit(ctx, func(ctx context.Context) error {
			// Simulate task work with a sleep.
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Task %d completed\n", index)
			return nil
		}, flexpool.MediumPriority)

		if err != nil {
			fmt.Printf("Failed to submit task %d: %v\n", index, err)
		}
	}

	// Wait for all submitted tasks to complete.
	pool.WaitUntilDone()

	// Optionally, handle errors reported by workers.
	go func() {
		for err := range pool.ErrorChan() {
			fmt.Printf("Error encountered: %v\n", err)
		}
	}()

	// Shutdown the pool gracefully, waiting for all workers to finish.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err := pool.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Printf("Failed to shutdown pool: %v\n", err)
	}

	fmt.Println("All tasks completed and pool shutdown gracefully.")
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

### Waiting for Tasks to Complete
To wait for all submitted tasks to complete:

```go
pool.WaitUntilDone()
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
go func() {
    for err := range pool.ErrorChan() {
        fmt.Printf("Error encountered: %v\n", err)
    }
}()
```

### Contributing
Contributions are welcome! Please submit issues or pull requests for any enhancements or bug fixes.

License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements

This project was inspired by various goroutine pool implementations and aims to provide a flexible and robust solution for managing concurrent tasks in Go.

---

Feel free to adjust any sections or add more details as necessary for your project.