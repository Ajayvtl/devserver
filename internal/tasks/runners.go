package tasks

import (
	"context"
)

// WorkspaceIndexRunner handles indexing tasks.
type WorkspaceIndexRunner struct{}

func (r *WorkspaceIndexRunner) Execute(ctx context.Context, task *Task, runtime *Runtime) error {
	workspaceID := task.WorkspaceID
	if workspaceID == "" {
		workspaceID = "devserver" // Default for now
	}

	runtime.Logger.Info().Str("task_id", task.ID).Msg("Starting workspace index runner")

	if runtime.Indexer != nil {
		if err := runtime.Indexer.Refresh(ctx, workspaceID); err != nil {
			runtime.Logger.Error().Err(err).Msg("Indexer refresh failed")
			return err
		}
	} else {
		runtime.Logger.Warn().Msg("Indexer is nil in runtime")
	}

	return nil
}

// WorkspaceSetupRunner handles setup tasks.
type WorkspaceSetupRunner struct{}

func (r *WorkspaceSetupRunner) Execute(ctx context.Context, task *Task, runtime *Runtime) error {
	runtime.Logger.Info().Str("task_id", task.ID).Msg("Starting workspace setup runner")
	return nil
}
