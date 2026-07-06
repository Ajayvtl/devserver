package capabilities

import (
	"context"
	"fmt"
	"sort"
)

// Resolver selects the best provider for a capability.
type Resolver struct {
	registry *Registry
}

func NewResolver(registry *Registry) *Resolver {
	return &Resolver{registry: registry}
}

func (r *Resolver) Resolve(cap Capability) ([]Binding, error) {
	if r == nil || r.registry == nil {
		return nil, fmt.Errorf("capability registry is unavailable")
	}
	bindings := r.registry.Describe(cap)
	sort.SliceStable(bindings, func(i, j int) bool {
		if bindings[i].Metadata.Priority == bindings[j].Metadata.Priority {
			return bindings[i].Metadata.Provider < bindings[j].Metadata.Provider
		}
		return bindings[i].Metadata.Priority > bindings[j].Metadata.Priority
	})
	return bindings, nil
}

func (r *Resolver) Select(ctx context.Context, cap Capability, required PermissionSet) (Binding, error) {
	bindings, err := r.Resolve(cap)
	if err != nil {
		return Binding{}, err
	}
	for _, binding := range bindings {
		if required != nil && !binding.Provider.Permissions(cap).AllowsSet(required) {
			continue
		}
		if binding.Provider == nil {
			continue
		}
		if err := binding.Provider.Health(ctx); err != nil {
			continue
		}
		return binding, nil
	}
	return Binding{}, fmt.Errorf("no provider found for capability %q", cap)
}
