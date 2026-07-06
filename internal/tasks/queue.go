package tasks

import (
	"container/heap"
	"sync"
)

// queue is a thread-safe priority queue for task definitions.
// FIFO is preserved within a priority level by sequence number.
type queue struct {
	mu        sync.Mutex
	items     priorityHeap
	maxSize   int
	closed    bool
	paused    bool
	seq       uint64
	notify    chan struct{}
	resume    chan struct{}
	closeOnce sync.Once
}

func newQueue(maxSize int) *queue {
	q := &queue{
		maxSize: maxSize,
		notify:  make(chan struct{}, 1),
		resume:  make(chan struct{}, 1),
	}
	heap.Init(&q.items)
	return q
}

// Enqueue adds a task to the queue. Returns false if full or closed.
func (q *queue) Enqueue(task *Task) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed || task == nil {
		return false
	}
	if q.maxSize > 0 && q.items.Len() >= q.maxSize {
		return false
	}
	q.seq++
	heap.Push(&q.items, &queuedItem{task: task, seq: q.seq})
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return true
}

// Dequeue removes and returns the highest-priority item, or nil if empty.
func (q *queue) Dequeue() *Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.paused || q.items.Len() == 0 {
		return nil
	}
	return heap.Pop(&q.items).(*queuedItem).task
}

// Len returns the number of items in the queue.
func (q *queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items.Len()
}

// Close prevents further enqueues and unblocks waiters.
func (q *queue) Close() {
	q.closeOnce.Do(func() {
		q.mu.Lock()
		q.closed = true
		q.mu.Unlock()
		close(q.notify)
		close(q.resume)
	})
}

// Pause suspends task draining.
func (q *queue) Pause() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.paused = true
}

// Resume re-enables task draining and wakes workers.
func (q *queue) Resume() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	wasPaused := q.paused
	q.paused = false
	q.mu.Unlock()
	if wasPaused {
		select {
		case q.resume <- struct{}{}:
		default:
		}
	}
}

// IsPaused reports whether the queue is currently paused.
func (q *queue) IsPaused() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.paused
}

// Remove cancels a queued task by ID. Returns true if found and removed.
func (q *queue) Remove(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if id == "" {
		return false
	}
	for i, item := range q.items {
		if item.task != nil && item.task.ID == id {
			heap.Remove(&q.items, i)
			return true
		}
	}
	return false
}

// Wait returns the notification channel. Consumers select on this to
// know when new items are available without busy-polling.
func (q *queue) Wait() <-chan struct{} {
	return q.notify
}

// ResumeWait notifies workers that paused draining may continue.
func (q *queue) ResumeWait() <-chan struct{} {
	return q.resume
}

// --- heap interface implementation ---

type queuedItem struct {
	task *Task
	seq uint64
}

type priorityHeap []*queuedItem

func (h priorityHeap) Len() int { return len(h) }
func (h priorityHeap) Less(i, j int) bool {
	if h[i].task.Priority == h[j].task.Priority {
		return h[i].seq < h[j].seq
	}
	return h[i].task.Priority < h[j].task.Priority
}
func (h priorityHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *priorityHeap) Push(x any) {
	*h = append(*h, x.(*queuedItem))
}

func (h *priorityHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return item
}
