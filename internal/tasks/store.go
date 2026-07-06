package tasks

import (
	"sync"
	"time"
)

// store is an in-memory task store, thread-safe.
// It holds the runtime state of all tasks known to the engine.
type store struct {
	mu      sync.RWMutex
	tasks   map[string]*Task
	order   []string
}

func newStore() *store {
	return &store{
		tasks: make(map[string]*Task),
	}
}

// Create adds a new task to the store.
func (s *store) Create(task *Task) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	if task.StartedAt.IsZero() {
		task.StartedAt = now
	}
	s.tasks[task.ID] = task
	s.order = append(s.order, task.ID)
	return task
}

// Get returns a task by ID, or nil if not found.
func (s *store) Get(id string) *Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

// Update persists changes to a task already in the store.
func (s *store) Update(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

// List returns all tasks in creation order.
func (s *store) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, 0, len(s.order))
	for _, id := range s.order {
		if t, ok := s.tasks[id]; ok {
			out = append(out, t)
		}
	}
	return out
}

// Delete removes a task by ID.
func (s *store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}
