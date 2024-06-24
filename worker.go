package flexpool

import (
	"context"
	"sync"
)

type Worker struct {
	id   int
	pool *Pool
	quit chan struct{}
	wg   sync.WaitGroup
}

func NewWorker(id int, pool *Pool) *Worker {
	return &Worker{
		id:   id,
		pool: pool,
		quit: make(chan struct{}),
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	defer w.wg.Done()

	for {
		select {
		case <-w.quit:
			return
		default:
			task := w.pool.getTask()
			if task == nil {
				// Pool is shutting down
				return
			}
			err := w.executeTask(task)
			if err != nil {
				w.handleError(err)
			}
		}
	}
}

func (w *Worker) executeTask(task *Task) error {
	return task.Execute()
}

func (w *Worker) handleError(err error) {
	// Log the error
	w.pool.logError(w.id, err)

	// You could also send the error to a channel for centralized handling
	// if w.pool.errorChannel != nil {
	//     select {
	//     case w.pool.errorChannel <- err:
	//     default:
	//         // Channel is full or closed, log locally
	//         log.Printf("Worker %d: Error channel full or closed. Error: %v", w.id, err)
	//     }
	// }
}

func (w *Worker) Stop(ctx context.Context) error {
	close(w.quit)

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
