package runtime

import (
	"context"
)

// Status represents the operational state of a Component.
type Status string

const (
	StatusStarting Status = "starting"
	StatusRunning  Status = "running"
	StatusStopped  Status = "stopped"
	StatusError    Status = "error"
)

// Health represents the health status of a Component.
type Health string

const (
	HealthHealthy   Health = "healthy"
	HealthDegraded  Health = "degraded"
	HealthUnhealthy Health = "unhealthy"
)

// Component is the minimal interface that any subsystem (Tasks, Providers, Workspaces)
// must implement to participate in the DevServer runtime lifecycle.
type Component interface {
	// Name returns the unique identifier for this component.
	Name() string

	// Initialize prepares the component (e.g., loading config, setting up structures).
	// It should not start long-running goroutines.
	Initialize(ctx context.Context) error

	// Start begins the component's main operation (e.g., listening on ports, starting workers).
	Start(ctx context.Context) error

	// Stop gracefully halts the component.
	Stop(ctx context.Context) error

	// Status returns the current operational state.
	Status() Status

	// Health returns the current health condition.
	Health() Health
}
