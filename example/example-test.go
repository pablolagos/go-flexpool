package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pablolagos/go-flexpool"
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

	// Shutdown the pool gracefully, waiting for all workers to finish.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err := pool.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Printf("Failed to shutdown pool: %v\n", err)
	}

	fmt.Println("All tasks completed and pool shutdown gracefully.")
}
