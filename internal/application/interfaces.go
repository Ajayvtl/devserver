package application

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/resource"
	"github.com/Ajayvtl/devserver/internal/domain/workflow"
)

// Orchestrator coordinates operations across repositories, executors, and workflows.
type Orchestrator interface {
	BootEnvironment(ctx context.Context, envID common.EnvironmentID) error
	SyncEnvironment(ctx context.Context, envID common.EnvironmentID) error
	DeployResource(ctx context.Context, spec *resource.ResourceSpec) error
	ExecuteWorkflow(ctx context.Context, wf *workflow.Workflow) error
}

// EventPublisher handles broadcasting domain events to external observers (e.g. UI).
type EventPublisher interface {
	Publish(ctx context.Context, eventID common.EventID) error
}
