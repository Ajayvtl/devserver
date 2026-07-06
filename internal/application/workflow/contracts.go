package workflow

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	domainWorkflow "github.com/Ajayvtl/devserver/internal/domain/workflow"
)

// Dispatcher is responsible for sending actions to the appropriate executors.
type Dispatcher interface {
	Dispatch(ctx context.Context, node *domainWorkflow.WorkflowNode) error
}

// Engine defines the complete lifecycle management for a workflow.
type Engine interface {
	CreateExecutionPlan(ctx context.Context, wf *domainWorkflow.Workflow) ([]*domainWorkflow.WorkflowNode, error)
	ExecutePlan(ctx context.Context, plan []*domainWorkflow.WorkflowNode) error
	Pause(ctx context.Context, wfID common.WorkflowID) error
	Resume(ctx context.Context, wfID common.WorkflowID) error
	Cancel(ctx context.Context, wfID common.WorkflowID) error
	Retry(ctx context.Context, wfID common.WorkflowID, nodeID string) error
	Rollback(ctx context.Context, wfID common.WorkflowID) error
	Status(ctx context.Context, wfID common.WorkflowID) (string, error)
}
