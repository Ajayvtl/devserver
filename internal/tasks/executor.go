package tasks

import "context"

// Executor runs a task definition.
type Executor interface {
	Run(ctx context.Context, def *Definition, tctx *Context) error
}

// directExecutor is a minimal fallback executor used when no
// core.Runner-backed executor has been injected.
type directExecutor struct{}

// NewDirectExecutor returns a basic executor that calls Fn directly.
func NewDirectExecutor() Executor {
	return directExecutor{}
}

func (directExecutor) Run(ctx context.Context, def *Definition, tctx *Context) error {
	if def == nil || def.Fn == nil {
		return nil
	}
	if err := def.Fn(tctx); err != nil {
		if def.RollbackFn == nil {
			return err
		}
		if rbErr := def.RollbackFn(tctx); rbErr != nil {
			return &ExecutionError{
				Err:        err,
				RolledBack: false,
			}
		}
		return &ExecutionError{Err: err, RolledBack: true}
	}
	return nil
}

// ExecutionError wraps a task failure and marks whether rollback succeeded.
type ExecutionError struct {
	Err        error
	RolledBack bool
}

func (e *ExecutionError) Error() string {
	if e == nil || e.Err == nil {
		return "task execution failed"
	}
	return e.Err.Error()
}

func (e *ExecutionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
