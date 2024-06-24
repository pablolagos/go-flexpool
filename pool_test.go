package flexpool_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pablolagos/go-flexpool"
)

func BenchmarkFlexPool(b *testing.B) {
	pool := flexpool.New(10, 1) // 10 workers, 100 maximum tasks
	b.ReportAllocs()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		err := pool.Submit(context.Background(), func(ctx context.Context) error {
			defer wg.Done()
			// Simulate task work
			time.Sleep(10 * time.Millisecond)
			return nil
		}, flexpool.MediumPriority)

		if err != nil {
			b.Errorf("Failed to submit task: %v", err)
		}
	}

	wg.Wait()
	pool.Shutdown(context.Background())
}

func main() {
	result := testing.Benchmark(BenchmarkFlexPool)
	fmt.Printf("%s\n", result)
}
