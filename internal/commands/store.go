package commands

import (
	"sync"
	"time"
)

type store struct {
	mu      sync.RWMutex
	records map[string]*Record
	order   []string
}

func newStore() *store {
	return &store{records: make(map[string]*Record)}
}

func (s *store) Create(cmd *Command) *Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	r := &Record{
		ID:          cmd.ID,
		Name:        cmd.Name,
		WorkspaceID: cmd.WorkspaceID,
		Capability:  cmd.Capability.String(),
		Provider:    cmd.Provider,
		Target:      cmd.Target,
		Status:      StatusQueued,
		Parameters:  cmd.Parameters,
		Metadata:    cmd.Metadata,
		CreatedAt:   now,
	}
	s.records[r.ID] = r
	s.order = append(s.order, r.ID)
	return r
}

func (s *store) Get(id string) *Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.records[id]
}

func (s *store) Update(r *Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[r.ID] = r
}

func (s *store) List() []*Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Record, 0, len(s.order))
	for _, id := range s.order {
		if record, ok := s.records[id]; ok {
			out = append(out, record)
		}
	}
	return out
}
