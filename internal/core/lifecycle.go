package core

import "context"

type Lifecycle struct {
	runner *Runner
}

func NewLifecycle(runner *Runner) *Lifecycle {
	return &Lifecycle{runner: runner}
}

func (l *Lifecycle) Execute(ctx context.Context, tasks ...Task) error {
	return l.runner.Run(ctx, tasks...)
}

func (l *Lifecycle) Stop(context.Context) error {
	return nil
}
