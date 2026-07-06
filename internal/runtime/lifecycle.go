package runtime

import (
	"context"
	"fmt"
)

// Coordinator manages the lifecycle phases of all registered components.
type Coordinator struct {
	registry *Registry
}

// NewCoordinator creates a new lifecycle coordinator.
func NewCoordinator(registry *Registry) *Coordinator {
	return &Coordinator{
		registry: registry,
	}
}

// Initialize calls Initialize on all components in registration order.
func (c *Coordinator) Initialize(ctx context.Context) error {
	for _, comp := range c.registry.All() {
		if err := comp.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize component %s: %w", comp.Name(), err)
		}
	}
	return nil
}

// Start calls Start on all components in registration order.
func (c *Coordinator) Start(ctx context.Context) error {
	components := c.registry.All()
	var started []Component

	for _, comp := range components {
		if err := comp.Start(ctx); err != nil {
			// Rollback previously started components
			for i := len(started) - 1; i >= 0; i-- {
				_ = started[i].Stop(context.Background())
			}
			return fmt.Errorf("failed to start component %s: %w", comp.Name(), err)
		}
		started = append(started, comp)
	}
	return nil
}

// Stop calls Stop on all components in reverse registration order.
func (c *Coordinator) Stop(ctx context.Context) error {
	components := c.registry.All()
	var lastErr error

	// Stop in reverse order
	for i := len(components) - 1; i >= 0; i-- {
		comp := components[i]
		if err := comp.Stop(ctx); err != nil {
			lastErr = fmt.Errorf("failed to stop component %s: %w", comp.Name(), err)
		}
	}
	return lastErr
}

type AggregatedHealth struct {
	Status     string            `json:"status"`
	Components map[string]string `json:"components"`
}

// AggregateHealth collects health status from all registered components.
func (c *Coordinator) AggregateHealth() AggregatedHealth {
	health := AggregatedHealth{
		Status:     "healthy",
		Components: make(map[string]string),
	}

	for _, comp := range c.registry.All() {
		compHealth := string(comp.Health())
		health.Components[comp.Name()] = compHealth
		if compHealth != string(HealthHealthy) {
			health.Status = "degraded"
		}
	}
	return health
}
