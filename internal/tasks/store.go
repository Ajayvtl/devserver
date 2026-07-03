package tasks

import (
	"sync"
	"time"
)

// store is an in-memory task record store, thread-safe.
// It holds the runtime state of all tasks known to the engine.
type store struct {
	mu      sync.RWMutex
	records map[string]*Record
	order   []string
}

func newStore() *store {
	return &store{
		records: make(map[string]*Record),
	}
}

// Create adds a new record from a definition.
func (s *store) Create(def *Definition) *Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	r := &Record{
		ID:          def.ID,
		Name:        def.Name,
		WorkspaceID: def.WorkspaceID,
		Status:      StatusQueued,
		Progress:    0,
		Detail:      "Queued",
		Priority:    def.Priority,
		Attempt:     0,
		MaxRetries:  def.Retry.MaxRetries,
		Metadata:    def.Metadata,
		CreatedAt:   now,
	}
	s.records[def.ID] = r
	s.order = append(s.order, def.ID)
	return r
}

// Get returns a record by ID, or nil if not found.
func (s *store) Get(id string) *Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.records[id]
}

// Update persists changes to a record already in the store.
func (s *store) Update(r *Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[r.ID] = r
}

// List returns all records in creation order.
func (s *store) List() []*Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Record, 0, len(s.order))
	for _, id := range s.order {
		if r, ok := s.records[id]; ok {
			out = append(out, r)
		}
	}
	return out
}

// Delete removes a record by ID.
func (s *store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}
