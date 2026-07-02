package core

import "github.com/Ajayvtl/devserver/internal/registry"

type Registry struct {
	inner *registry.Registry
}

func NewRegistry() *Registry {
	return &Registry{inner: registry.New()}
}

func (r *Registry) Register(module registry.Module) error {
	return r.inner.Register(module)
}

func (r *Registry) Lookup(name string) (registry.Module, bool) {
	return r.inner.Lookup(name)
}

func (r *Registry) Modules() []registry.Module {
	return r.inner.Modules()
}

func (r *Registry) Inner() *registry.Registry {
	return r.inner
}
