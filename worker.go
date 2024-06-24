package flexpool

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
)

type Worker struct {
	id     int
	pool   *Pool
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorker(id int, pool *Pool) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		id:     id,
		pool:   pool,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				return
			default:
				w.pool.mutex.Lock()
				for w.pool.tasks.Len() == 0 && atomic.LoadInt32(&w.pool.running) == 1 {
					w.pool.cond.Wait()
				}
				if atomic.LoadInt32(&w.pool.running) == 0 {
					w.pool.mutex.Unlock()
					return
				}
				task := heap.Pop(w.pool.tasks).(*Task)
				w.pool.mutex.Unlock()

				if err := task.Execute(); err != nil {
					// Log or handle the error
				}
			}
		}
	}()
}

func (w *Worker) Stop(ctx context.Context) {
	w.cancel()
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}
