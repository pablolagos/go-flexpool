package flexpool_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pablolagos/go-flexpool"
)

func BenchmarkFlexPool(b *testing.B) {
	pool := flexpool.New(10, 0) // 10 workers, 100 maximum tasks

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := pool.Submit(context.Background(), func(ctx context.Context) error {
			// Simulate task work
			time.Sleep(10 * time.Millisecond)
			return nil
		}, flexpool.MediumPriority)

		if err != nil {
			b.Errorf("Failed to submit task: %v", err)
		}
	}

	// Wait for all tasks to complete
	pool.WaitUntilDone()

	// Shutdown the pool
	pool.Shutdown(context.Background())
}

func main() {
	result := testing.Benchmark(BenchmarkFlexPool)
	fmt.Printf("%s\n", result)
}
