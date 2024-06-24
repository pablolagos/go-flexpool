package flexpool

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"sync/atomic"
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
	maxWorkers int32
	maxTasks   int32
	running    int32
}

type Task struct {
	Execute  func() error
	Priority Priority
}

func New(maxWorkers, maxTasks int) *Pool {
	p := &Pool{
		maxWorkers: int32(maxWorkers),
		maxTasks:   int32(maxTasks),
		tasks:      &PriorityQueue{},
		running:    1,
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

func (p *Pool) Submit(ctx context.Context, task func() error, priority Priority) error {
	if atomic.LoadInt32(&p.running) == 0 {
		return ErrPoolClosed
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.tasks.Len() >= int(p.maxTasks) {
		return ErrPoolFull
	}

	t := &Task{Execute: task, Priority: priority}
	heap.Push(p.tasks, t)
	p.cond.Signal()

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
			w := NewWorker(int(i), p)
			p.workers = append(p.workers, w)
			go w.Start()
		}
	} else if int32(newSize) < oldSize {
		// Remove excess workers
		for i := newSize; i < int(oldSize); i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				p.workers[i].Stop(ctx)
			}
		}
		p.workers = p.workers[:newSize]
	}

	atomic.StoreInt32(&p.maxWorkers, int32(newSize))
	return nil
}

func (p *Pool) SetMaxTasks(newMax int) {
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
