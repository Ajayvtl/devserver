package providers

import (
	"context"
	"fmt"
	"sync"

	rt "github.com/Ajayvtl/devserver/internal/runtime"
)

// Manager manages a collection of registered providers.
type Manager struct {
	mu        sync.RWMutex
	providers map[string]Provider
	status    rt.Status
}

// NewManager creates a new Provider Manager.
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
		status:    rt.StatusStopped,
	}
}

func (m *Manager) Name() string { return "providers.Manager" }

func (m *Manager) Initialize(ctx context.Context) error {
	m.status = rt.StatusStarting
	return nil
}

func (m *Manager) Start(ctx context.Context) error {
	m.status = rt.StatusRunning
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	m.status = rt.StatusStopped
	return nil
}

func (m *Manager) Status() rt.Status { return m.status }

func (m *Manager) Health() rt.Health { return rt.HealthHealthy }

// Register adds a provider to the manager.
func (m *Manager) Register(p Provider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[p.Metadata().Name] = p
}

// Get retrieves a provider by name.
func (m *Manager) Get(name string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if p, ok := m.providers[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider %q not found", name)
}

// List returns all registered providers.
func (m *Manager) List() []Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []Provider
	for _, p := range m.providers {
		list = append(list, p)
	}
	return list
}
