package flexpool

// PriorityQueue is a heap-based priority queue of tasks.
type PriorityQueue []*Task

// Len returns the number of tasks in the queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares the priority of two tasks.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

// Swap exchanges the positions of two tasks in the queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds a task to the queue.
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}

// Pop removes and returns the highest priority task from the queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
