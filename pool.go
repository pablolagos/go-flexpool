package flexpool

import (
	"container/heap"
	"context"
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
	mutex      sync.RWMutex
	cond       *sync.Cond
	maxWorkers int32
	maxTasks   int32
	running    int32
	errorChan  chan error // Optional: channel for error reporting
}

type Task struct {
	Execute  func(ctx context.Context) error
	Priority Priority
	ctx      context.Context
}

func New(maxWorkers, maxTasks int) *Pool {
	p := &Pool{
		maxWorkers: int32(maxWorkers),
		maxTasks:   int32(maxTasks),
		tasks:      &PriorityQueue{},
		running:    1,
		errorChan:  make(chan error, 100), // Buffer size can be adjusted
	}
	p.cond = sync.NewCond(&p.mutex)
	heap.Init(p.tasks)

	for i := 0; i < maxWorkers; i++ {
		w := NewWorker(i, p)
		p.workers = append(p.workers, w)
		go w.Start()
	}

	return p
}

func (p *Pool) Submit(ctx context.Context, task func(ctx context.Context) error, priority Priority) error {
	if atomic.LoadInt32(&p.running) == 0 {
		return ErrPoolClosed
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	maxTasks := atomic.LoadInt32(&p.maxTasks)
	if maxTasks != UnlimitedTasks && p.tasks.Len() >= int(maxTasks) {
		return ErrPoolFull
	}

	t := &Task{Execute: task, Priority: priority, ctx: ctx}
	heap.Push(p.tasks, t)
	p.cond.Signal()

	return nil
}

func (p *Pool) getTask() *Task {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for p.tasks.Len() == 0 && atomic.LoadInt32(&p.running) == 1 {
		p.cond.Wait()
	}

	if atomic.LoadInt32(&p.running) == 0 {
		return nil
	}

	if p.tasks.Len() > 0 {
		return heap.Pop(p.tasks).(*Task)
	}

	return nil
}

func (p *Pool) Resize(ctx context.Context, newSize int) error {
	if newSize < 0 {
		return errors.New("new size must be non-negative")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	oldSize := atomic.LoadInt32(&p.maxWorkers)
	if int32(newSize) > oldSize {
		// Add more workers
		for i := oldSize; i < int32(newSize); i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				w := NewWorker(int(i), p)
				p.workers = append(p.workers, w)
				go w.Start()
			}
		}
	} else if int32(newSize) < oldSize {
		// Remove excess workers
		var wg sync.WaitGroup
		for i := newSize; i < int(oldSize); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				p.workers[index].Stop(ctx)
			}(i)
		}
		wg.Wait()
		p.workers = p.workers[:newSize]
	}

	atomic.StoreInt32(&p.maxWorkers, int32(newSize))
	return nil
}

func (p *Pool) SetMaxTasks(newMax int) {
	if newMax < 0 {
		newMax = 0 // Convert negative values to 0 (unlimited)
	}
	atomic.StoreInt32(&p.maxTasks, int32(newMax))
}

func (p *Pool) Shutdown(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
		return ErrPoolClosed
	}

	p.cond.Broadcast()

	var wg sync.WaitGroup
	for _, worker := range p.workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.Stop(ctx)
		}(worker)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

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

func (p *Pool) QueuedTasks() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.tasks.Len()
}
