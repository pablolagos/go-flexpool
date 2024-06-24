package main

import "container/heap"

type Worker struct {
	id   int
	pool *Pool
	quit chan struct{}
}

func (w *Worker) Start() {
	w.quit = make(chan struct{})
	for {
		w.pool.mutex.Lock()
		for w.pool.tasks.Len() == 0 && w.pool.running {
			w.pool.cond.Wait()
		}
		if !w.pool.running {
			w.pool.mutex.Unlock()
			return
		}
		task := heap.Pop(w.pool.tasks).(*Task)
		w.pool.mutex.Unlock()

		select {
		case <-w.quit:
			return
		default:
			task.Execute()
		}
	}
}

func (w *Worker) Stop() {
	close(w.quit)
}
