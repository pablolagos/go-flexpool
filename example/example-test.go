package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	pool := New(5, 20)

	// Submit tasks
	for i := 0; i < 10; i++ {
		i := i
		priority := Priority(i % 3)
		err := pool.Submit(func() error {
			fmt.Printf("Task %d (Priority %d) started\n", i, priority)
			time.Sleep(time.Second)
			fmt.Printf("Task %d (Priority %d) completed\n", i, priority)
			return nil
		}, priority)
		if err != nil {
			fmt.Printf("Failed to submit task: %v\n", err)
		}
	}

	// Resize the pool
	pool.Resize(8)

	// Change max tasks
	pool.SetMaxTasks(30)

	// Submit more tasks...

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool.Shutdown(ctx)
}
