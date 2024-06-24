// pool.go
package workerpool

import (
	"container/heap"
	"context"
	"errors"
	"sync"
)

var (
	ErrPoolFull   = errors.New("task queue is full")
	ErrPoolClosed = errors.New("pool is closed")
)

type Priority int

const (
	LowPriority Priority = iota
	MediumPriority
	HighPriority
)

type Pool struct {
	workers    []*Worker
	tasks      *PriorityQueue
	mutex      sync.RWMutex
	cond       *sync.Cond
	quit       chan struct{}
	maxWorkers int
	maxTasks   int
	running    bool
}

type Task struct {
	Execute  func() error
	Priority Priority
}

func New(maxWorkers, maxTasks int) *Pool {
	p := &Pool{
		maxWorkers: maxWorkers,
		maxTasks:   maxTasks,
		tasks:      &PriorityQueue{},
		quit:       make(chan struct{}),
		running:    true,
	}
	p.cond = sync.NewCond(&p.mutex)
	heap.Init(p.tasks)

	for i := 0; i < maxWorkers; i++ {
		w := &Worker{
			id:   i,
			pool: p,
		}
		p.workers = append(p.workers, w)
		go w.Start()
	}

	return p
}

func (p *Pool) Submit(task func() error, priority Priority) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.running {
		return ErrPoolClosed
	}

	if p.tasks.Len() >= p.maxTasks {
		return ErrPoolFull
	}

	heap.Push(p.tasks, &Task{Execute: task, Priority: priority})
	p.cond.Signal()
	return nil
}

func (p *Pool) Shutdown(ctx context.Context) {
	p.mutex.Lock()
	p.running = false
	p.mutex.Unlock()

	close(p.quit)

	done := make(chan struct{})
	go func() {
		p.cond.Broadcast()
		for _, worker := range p.workers {
			worker.Stop()
		}
		close(done)
	}()

	select {
	case <-done:
		// All workers finished
	case <-ctx.Done():
		// Timeout reached
	}
}

func (p *Pool) Resize(newSize int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if newSize > p.maxWorkers {
		// Add more workers
		for i := p.maxWorkers; i < newSize; i++ {
			w := &Worker{
				id:   i,
				pool: p,
			}
			p.workers = append(p.workers, w)
			go w.Start()
		}
	} else if newSize < p.maxWorkers {
		// Remove excess workers
		for i := newSize; i < p.maxWorkers; i++ {
			p.workers[i].Stop()
		}
		p.workers = p.workers[:newSize]
	}
	p.maxWorkers = newSize
}

func (p *Pool) SetMaxTasks(newMax int) {
	p.mutex.Lock()
	p.maxTasks = newMax
	p.mutex.Unlock()
}

// worker.go
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

// priority_queue.go
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
