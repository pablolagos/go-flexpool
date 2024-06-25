# go-flexpool

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pablolagos/go-flexpool)](https://goreportcard.com/report/github.com/pablolagos/go-flexpool)
[![GoDoc](https://godoc.org/github.com/pablolagos/go-flexpool?status.svg)](https://godoc.org/github.com/pablolagos/go-flexpool)

`go-flexpool` is a flexible and efficient worker pool implementation in Go, designed to handle tasks with different priorities seamlessly.

## Features

- **Priority-Based Task Scheduling**: Submit tasks with different priorities (High, Medium, Low).
- **Dynamic Worker Pool Resizing**: Increase or decrease the number of workers on-the-fly.
- **Graceful Shutdown**: Ensure all tasks are completed before shutting down the pool.
- **Error Handling**: Access errors encountered during task execution through an error channel.
- **Concurrency**: Efficiently manage concurrent task execution with Go's goroutines and channels.

## Installation

Install `go-flexpool` using `go get`:

```go
go get github.com/pablolagos/go-flexpool
```

## Usage

### Create a New Pool

Create a new pool with a specified number of workers and maximum tasks:

```go
import "github.com/pablolagos/go-flexpool"

pool := flexpool.New(100, 1000) // 100 workers, max 1000 queued tasks
```

### Submit Tasks

Submit tasks to the pool with different priorities:

```go
err := pool.Submit(func() error {
    // Your task code here
    return nil
}, flexpool.MediumPriority)

if err != nil {
log.Printf("Failed to submit task: %v", err)
}
```

### Resize the Pool

Resize the pool to adjust the number of workers dynamically:

```go
err := pool.Resize(200) // Resize to 200 workers
if err != nil {
    log.Printf("Failed to resize pool: %v", err)
}
```

### Shutdown the Pool

Gracefully shutdown the pool, ensuring all tasks are completed:

```go
err := pool.Shutdown()
if err != nil {
    log.Printf("Failed to shutdown pool: %v", err)
}
```

### Wait Until All Tasks are Done

Wait until all submitted tasks are completed:

```go
pool.WaitUntilDone()
```

### Access Errors

Access errors encountered during task execution through the error channel:

```go
for err := range pool.Errors() {
    log.Printf("Task error: %v", err)
}
```

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or suggestions, feel free to open an issue or contact me at [pablolagos](https://github.com/pablolagos).

---

<p align="center">
  <a href="https://github.com/pablolagos/go-flexpool">GitHub Repository</a>
</p>