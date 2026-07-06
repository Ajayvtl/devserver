package contracts

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
)

// Executor defines the canonical interface for all execution targets.
type Executor interface {
	// Execute performs the requested action against the environment.
	Execute(ctx context.Context, act *action.Action) (*ExecutionResult, error)

	// Detect queries the environment to see if the executor can operate.
	Detect(ctx context.Context) (bool, []Diagnostic, error)

	// Capabilities returns the typed capabilities supported by this executor.
	Capabilities() Capabilities

	// Validate ensures the action can be run by this executor before execution.
	Validate(ctx context.Context, act *action.Action) error

	// Health checks if the executor is still alive and responsive.
	Health(ctx context.Context) (common.HealthState, error)

	// Metadata returns descriptive information about this executor instance.
	Metadata() ExecutorMetadata
}

// ExecutorMetadata describes a specific instance of an executor.
type ExecutorMetadata struct {
	ID      string
	Type    common.ExecutorType
	Version string
}
