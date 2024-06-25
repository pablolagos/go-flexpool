package flexpool

import (
	"sync"
)

type Worker struct {
	id   int
	pool *Pool
	quit chan struct{}
	wg   sync.WaitGroup
}

// NewWorker creates a new worker associated with the given pool.
func NewWorker(id int, pool *Pool) *Worker {
	return &Worker{
		id:   id,
		pool: pool,
		quit: make(chan struct{}),
	}
}

// Start begins the worker's task processing loop.
func (w *Worker) Start() {
	w.wg.Add(1)
	defer w.wg.Done()

	for {
		select {
		case <-w.quit:
			return
		case task, ok := <-w.pool.taskChan:
			if !ok {
				// Pool is shutting down
				return
			}
			err := w.executeTask(task)
			if err != nil {
				w.handleError(err)
			}
			w.pool.waitGroup.Done()
		}
	}
}

// executeTask executes the given task and returns any error encountered.
func (w *Worker) executeTask(task *Task) error {
	return task.Execute()
}

// handleError logs an error encountered during task execution.
func (w *Worker) handleError(err error) {
	// Log the error
	w.pool.logError(w.id, err)
}

// Stop signals the worker to stop processing tasks and waits for it to finish.
func (w *Worker) Stop() {
	close(w.quit)
	w.wg.Wait()
}
