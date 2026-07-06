package contracts

import (
	"context"
	"errors"
	"sync"

	"github.com/Ajayvtl/devserver/internal/application/envcontext"
)

var ErrExecutorNotFound = errors.New("executor not found in registry")

// ExecutorRegistry manages available executors.
type ExecutorRegistry interface {
	Register(ctx context.Context, exec Executor) error
	Get(ctx context.Context, id string) (Executor, error)
	List(ctx context.Context) ([]Executor, error)
	Unregister(ctx context.Context, id string) error
}

type memoryRegistry struct {
	mu        sync.RWMutex
	executors map[string]Executor
}

func NewMemoryRegistry() ExecutorRegistry {
	return &memoryRegistry{
		executors: make(map[string]Executor),
	}
}

func (r *memoryRegistry) Register(ctx context.Context, exec Executor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executors[exec.Metadata().ID] = exec
	return nil
}

func (r *memoryRegistry) Get(ctx context.Context, id string) (Executor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exec, exists := r.executors[id]
	if !exists {
		return nil, ErrExecutorNotFound
	}
	return exec, nil
}

func (r *memoryRegistry) List(ctx context.Context) ([]Executor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []Executor
	// If context has an active environment, we use it to filter executors.
	// We create a dummy switcher to just extract the context key.
	// Actually, we can just use a generic func or the switcher.
	// A cleaner way is just instantiate envcontext.NewSwitcher(nil) to call FromContext.
	switcher := envcontext.NewSwitcher(nil)
	targetEnv, hasEnv := switcher.FromContext(ctx)

	for _, exec := range r.executors {
		if hasEnv {
			if exec.Metadata().EnvironmentID != targetEnv {
				continue
			}
		}
		list = append(list, exec)
	}
	return list, nil
}

func (r *memoryRegistry) Unregister(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.executors, id)
	return nil
}
