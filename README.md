# go-flexpool

A flexible, priority-based worker pool for Go with dynamic resizing capabilities.

## Why Another Worker Pool?

While there are several worker pool implementations available for Go, go-flexpool stands out for the following reasons:

1. **Priority-based execution**: Tasks can be assigned priorities, ensuring that critical work is completed first.
2. **Dynamic resizing**: The number of workers can be adjusted at runtime, allowing the pool to adapt to changing workloads.
3. **Adjustable queue size**: The maximum number of queued tasks can be modified on-the-fly, providing control over memory usage.
4. **Simplicity**: The API is straightforward and easy to use, requiring minimal setup.
5. **Efficiency**: Utilizes Go's concurrency primitives for optimal performance.
6. **No external dependencies**: Built using only the Go standard library, minimizing potential conflicts and security risks.

go-flexpool aims to provide a balance between functionality and simplicity, offering advanced features without sacrificing ease of use or introducing unnecessary complexity.

## Features

- Priority-based task execution
- Dynamic pool size adjustment
- Adjustable task queue size
- Simple and efficient API
- No external dependencies

## Installation

```bash
go get github.com/yourusername/go-flexpool
```

## Quick Start

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/yourusername/go-flexpool"
)

func main() {
    // Create a new pool with 5 workers and a maximum of 20 tasks
    pool := flexpool.New(5, 20)

    // Submit a task with medium priority
    err := pool.Submit(func() error {
        fmt.Println("Executing task")
        return nil
    }, flexpool.MediumPriority)

    if err != nil {
        fmt.Printf("Failed to submit task: %v\n", err)
    }

    // Resize the pool to 8 workers
    pool.Resize(8)

    // Change max tasks to 30
    pool.SetMaxTasks(30)

    // Shutdown the pool with a 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    pool.Shutdown(ctx)
}
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.
