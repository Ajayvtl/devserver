package runtime

import (
	"fmt"
	"sync"
)

// Registry manages the discovery and registration of Components.
type Registry struct {
	mu         sync.RWMutex
	components map[string]Component
	order      []string
}

// NewRegistry creates a new component registry.
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[string]Component),
		order:      make([]string, 0),
	}
}

// Register adds a component to the registry.
func (r *Registry) Register(c Component) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := c.Name()
	if _, exists := r.components[name]; exists {
		return fmt.Errorf("component already registered: %s", name)
	}

	r.components[name] = c
	r.order = append(r.order, name)
	return nil
}

// Get retrieves a component by name.
func (r *Registry) Get(name string) (Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.components[name]
	if !exists {
		return nil, fmt.Errorf("component not found: %s", name)
	}
	return c, nil
}

// All returns all registered components in the order they were registered.
func (r *Registry) All() []Component {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []Component
	for _, name := range r.order {
		list = append(list, r.components[name])
	}
	return list
}
