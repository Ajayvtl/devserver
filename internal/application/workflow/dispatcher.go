package workflow

import (
	"context"
	"fmt"

	domainWorkflow "github.com/Ajayvtl/devserver/internal/domain/workflow"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	"github.com/Ajayvtl/devserver/internal/repository"
)

type defaultDispatcher struct {
	actions  repository.ActionRepository
	registry contracts.ExecutorRegistry
}

func NewDispatcher(actions repository.ActionRepository, registry contracts.ExecutorRegistry) Dispatcher {
	return &defaultDispatcher{
		actions:  actions,
		registry: registry,
	}
}

func (d *defaultDispatcher) Dispatch(ctx context.Context, node *domainWorkflow.WorkflowNode) error {
	// Look up action
	act, err := d.actions.GetByID(ctx, node.ActionID)
	if err != nil {
		return fmt.Errorf("failed to get action: %w", err)
	}

	// Try to execute via Command Bus if it's a known capability
	// Or use Executor Registry
	execs, err := d.registry.List(ctx)
	if err != nil {
		return err
	}

	for _, exec := range execs {
		if err := exec.Validate(ctx, act); err == nil {
			_, execErr := exec.Execute(ctx, act)
			return execErr
		}
	}

	return fmt.Errorf("no suitable executor found for action %s", act.ID)
}
