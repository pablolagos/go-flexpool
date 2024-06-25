package flexpool

import (
	"container/heap"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

var (
	ErrPoolFull   = errors.New("task queue is full")
	ErrPoolClosed = errors.New("pool is closed")
)

const UnlimitedTasks int32 = 0

type Priority int

const (
	LowPriority Priority = iota
	MediumPriority
	HighPriority
)

type Pool struct {
	workers    []*Worker
	tasks      *PriorityQueue
	taskChan   chan *Task
	maxWorkers int32
	maxTasks   int32
	running    int32
	errorChan  chan error // Optional: channel for error reporting
	waitGroup  sync.WaitGroup
}

type Task struct {
	Execute  func() error
	Priority Priority
}

// New creates a new worker pool with a specified number of workers and maximum tasks.
func New(maxWorkers, maxTasks int) *Pool {
	taskChan := make(chan *Task, maxTasks)
	p := &Pool{
		maxWorkers: int32(maxWorkers),
		maxTasks:   int32(maxTasks),
		tasks:      &PriorityQueue{},
		taskChan:   taskChan,
		running:    1,
		errorChan:  make(chan error, 100), // Buffer size can be adjusted
	}
	heap.Init(p.tasks)

	for i := 0; i < maxWorkers; i++ {
		w := NewWorker(i, p)
		p.workers = append(p.workers, w)
		go w.Start()
	}

	return p
}

// Submit adds a task to the pool with a specified priority. If the pool is full, it returns ErrPoolFull.
// If the pool is closed, it returns ErrPoolClosed.
func (p *Pool) Submit(task func() error, priority Priority) error {
	if atomic.LoadInt32(&p.running) == 0 {
		return ErrPoolClosed
	}

	p.waitGroup.Add(1)

	t := &Task{Execute: task, Priority: priority}
	select {
	case p.taskChan <- t:
		return nil
	default:
		p.waitGroup.Done()
		return ErrPoolFull
	}
}

// Resize changes the number of workers in the pool. It adds or removes workers as needed.
func (p *Pool) Resize(newSize int) error {
	if newSize < 0 {
		return errors.New("new size must be non-negative")
	}

	oldSize := atomic.LoadInt32(&p.maxWorkers)
	if int32(newSize) > oldSize {
		// Add more workers
		for i := oldSize; i < int32(newSize); i++ {
			w := NewWorker(int(i), p)
			p.workers = append(p.workers, w)
			go w.Start()
		}
	} else if int32(newSize) < oldSize {
		// Remove excess workers
		for i := newSize; i < int(oldSize); i++ {
			p.workers[i].Stop()
		}
		p.workers = p.workers[:newSize]
	}

	atomic.StoreInt32(&p.maxWorkers, int32(newSize))
	return nil
}

// SetMaxTasks sets the maximum number of tasks that can be queued. 0 means unlimited tasks.
func (p *Pool) SetMaxTasks(newMax int) {
	if newMax < 0 {
		newMax = 0 // Convert negative values to 0 (unlimited)
	}
	atomic.StoreInt32(&p.maxTasks, int32(newMax))
}

func (p *Pool) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
		return ErrPoolClosed
	}

	close(p.taskChan)

	var wg sync.WaitGroup
	for _, worker := range p.workers {
		wg.Add(1)
		go func(w *Worker) {
			w.Stop()
			wg.Done()
		}(worker)
	}

	wg.Wait()
	close(p.errorChan)
	return nil
}

// WaitUntilDone blocks until all submitted tasks have been completed.
func (p *Pool) WaitUntilDone() {
	p.waitGroup.Wait()
}

// Errors returns a read-only channel for receiving errors from tasks.
func (p *Pool) Errors() <-chan error {
	return p.errorChan
}

// logError logs an error encountered by a worker and optionally sends it to the error channel.
func (p *Pool) logError(workerID int, err error) {
	log.Printf("Worker %d encountered an error: %v", workerID, err)
	// Optionally, send to error channel
	select {
	case p.errorChan <- err:
	default:
		// Channel is full or closed, log locally
		log.Printf("Error channel full or closed. Error from Worker %d: %v", workerID, err)
	}
}

// QueuedTasks returns the number of tasks currently queued.
func (p *Pool) QueuedTasks() int {
	return len(p.taskChan)
}
