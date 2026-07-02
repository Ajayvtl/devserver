package state

import (
	"context"
	"sync"
	"time"
)

type Installation struct {
	Version   string            `yaml:"version" json:"version"`
	UpdatedAt time.Time         `yaml:"updated_at" json:"updated_at"`
	Metadata  map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Status    string            `yaml:"status" json:"status"`
}

type State struct {
	Installed map[string]Installation `yaml:"installed" json:"installed"`
}

type TaskEvent struct {
	TaskID   string `json:"taskId"`
	Name     string `json:"name"`
	Progress int    `json:"progress"`
	State    string `json:"state"`
	Detail   string `json:"detail"`
	Scope    string `json:"scope"`
}

func New() *State {
	return &State{Installed: map[string]Installation{}}
}

func (s *State) MarkInstalled(name, version string) {
	if s.Installed == nil {
		s.Installed = map[string]Installation{}
	}

	s.Installed[name] = Installation{
		Version:   version,
		UpdatedAt: time.Now().UTC(),
		Status:    "installed",
	}
}

type Store interface {
	Load(context.Context) (*State, error)
	Save(context.Context, *State) error
}

type MemoryStore struct {
	mu    sync.Mutex
	state *State
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{state: New()}
}

func (s *MemoryStore) Load(context.Context) (*State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return clone(s.state), nil
}

func (s *MemoryStore) Save(_ context.Context, st *State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = clone(st)
	return nil
}

func clone(st *State) *State {
	if st == nil {
		return New()
	}

	out := New()
	for name, installation := range st.Installed {
		out.Installed[name] = installation
	}
	return out
}
