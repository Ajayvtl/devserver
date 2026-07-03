package core

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

type runnerTaskExecutor struct {
	runner *Runner
}

// NewTaskExecutor returns a task executor that reuses the shared Runner.
func NewTaskExecutor(log zerolog.Logger) tasks.Executor {
	return &runnerTaskExecutor{runner: NewRunner(log)}
}

func (e *runnerTaskExecutor) Run(ctx context.Context, def *tasks.Definition, tctx *tasks.Context) error {
	if e == nil || e.runner == nil {
		return tasks.NewDirectExecutor().Run(ctx, def, tctx)
	}
	return e.runner.Run(ctx, runnerAdapter{def: def, tctx: tctx})
}

type runnerAdapter struct {
	def  *tasks.Definition
	tctx *tasks.Context
}

func (a runnerAdapter) Name() string {
	if a.def == nil {
		return "task"
	}
	return a.def.Name
}

func (a runnerAdapter) Run(ctx context.Context) error {
	if a.def == nil || a.def.Fn == nil {
		return nil
	}
	if err := a.def.Fn(a.tctx); err != nil {
		if a.def.RollbackFn == nil {
			return err
		}
		if rbErr := a.def.RollbackFn(a.tctx); rbErr != nil {
			return &tasks.ExecutionError{Err: err, RolledBack: false}
		}
		return &tasks.ExecutionError{Err: err, RolledBack: true}
	}
	return nil
}

func (a runnerAdapter) Rollback(ctx context.Context) error {
	if a.def == nil || a.def.RollbackFn == nil {
		return nil
	}
	return a.def.RollbackFn(a.tctx)
}
