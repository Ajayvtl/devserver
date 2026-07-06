package capabilities

import (
	"context"
	"fmt"
	"sync"
)

// Provider exposes a capability surface.
type Provider interface {
	Name() string
	Capabilities() []Capability
	Metadata(Capability) Metadata
	Permissions(Capability) PermissionSet
	Execute(context.Context, Capability, any) error
	Health(context.Context) error
}

// Registry indexes providers by capability.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
	byCap     map[Capability][]string
	order     []string
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
		byCap:     make(map[Capability][]string),
		order:     []string{},
	}
}

func (r *Registry) Register(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	name := provider.Name()
	if name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}

	r.providers[name] = provider
	r.order = append(r.order, name)
	for _, cap := range provider.Capabilities() {
		r.byCap[cap] = append(r.byCap[cap], name)
	}
	return nil
}

func (r *Registry) Lookup(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	return provider, ok
}

func (r *Registry) Providers() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Provider, 0, len(r.order))
	for _, name := range r.order {
		if provider, ok := r.providers[name]; ok {
			out = append(out, provider)
		}
	}
	return out
}

func (r *Registry) Resolve(cap Capability) []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := r.byCap[cap]
	out := make([]Provider, 0, len(names))
	for _, name := range names {
		if provider, ok := r.providers[name]; ok {
			out = append(out, provider)
		}
	}
	return out
}

func (r *Registry) Describe(cap Capability) []Binding {
	providers := r.Resolve(cap)
	out := make([]Binding, 0, len(providers))
	for _, provider := range providers {
		meta := provider.Metadata(cap)
		meta.Provider = provider.Name()
		meta.Capability = cap
		out = append(out, Binding{Provider: provider, Metadata: meta})
	}
	return out
}
