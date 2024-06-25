package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pablolagos/go-flexpool"
)

func main() {
	// Create a new pool with 100 workers and a maximum of 1000 queued tasks
	pool := flexpool.New(100, 1000)

	// Define a sample task
	task := func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("Task completed")
		return nil
	}

	// Submit the task with medium priority
	for i := 0; i < 10; i++ { // Submit multiple tasks
		err := pool.Submit(task, flexpool.MediumPriority)
		if err != nil {
			log.Printf("Failed to submit task: %v", err)
		}
	}

	// Wait until all submitted tasks are done
	pool.WaitUntilDone()

	// Shutdown the pool
	err := pool.Shutdown()
	if err != nil {
		log.Printf("Failed to shutdown pool: %v", err)
	}

	// Handle errors (if any)
	for err := range pool.Errors() {
		log.Printf("Task error: %v", err)
	}
}
