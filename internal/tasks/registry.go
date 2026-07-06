package tasks

import "sync"

// Registry stores the registered runners that can execute tasks.
type Registry struct {
	mu      sync.RWMutex
	runners map[string]Runner
}

// NewRegistry creates a new empty Runner registry.
func NewRegistry() *Registry {
	return &Registry{
		runners: make(map[string]Runner),
	}
}

// Register adds a runner to the registry for a specific task type.
func (r *Registry) Register(taskType string, runner Runner) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runners[taskType] = runner
}

// Runner returns the registered runner for the given task type, or false if not found.
func (r *Registry) Runner(taskType string) (Runner, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	runner, ok := r.runners[taskType]
	return runner, ok
}
