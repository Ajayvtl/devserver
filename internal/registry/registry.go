package registry

import (
	"context"
	"fmt"
	"sync"
)

type Module interface {
	Name() string
	Description() string
	Check(context.Context) error
	Install(context.Context) error
	Configure(context.Context) error
	Validate(context.Context) error
	Upgrade(context.Context) error
	Uninstall(context.Context) error
	Rollback(context.Context) error
}

type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
	order   []string
}

func New() *Registry {
	return &Registry{
		modules: make(map[string]Module),
		order:   []string{},
	}
}

func (r *Registry) Register(module Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := module.Name()
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q already registered", name)
	}

	r.modules[name] = module
	r.order = append(r.order, name)
	return nil
}

func (r *Registry) Lookup(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	module, ok := r.modules[name]
	return module, ok
}

func (r *Registry) Modules() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Module, 0, len(r.order))
	for _, name := range r.order {
		if module, ok := r.modules[name]; ok {
			out = append(out, module)
		}
	}
	return out
}
