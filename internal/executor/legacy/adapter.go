package legacy

import (
	"context"
)

type localAdapter struct{}

func (l *localAdapter) Run(task interface{}) error {
	return nil
}
func (l *localAdapter) Execute(ctx context.Context, cmd string, args ...string) (*LegacyExecutionResult, error) {
	return &LegacyExecutionResult{}, nil
}
func (l *localAdapter) StartProcess(ctx context.Context, cmd string, args ...string) error {
	return nil
}
func (l *localAdapter) StopProcess(ctx context.Context, name string) error {
	return nil
}
func (l *localAdapter) CheckHealth(ctx context.Context, name string) bool {
	return true
}
func (l *localAdapter) DetectBinary(ctx context.Context, name string) (bool, error) {
	return true, nil
}
func (l *localAdapter) QueryVersion(ctx context.Context, cmd string, args ...string) (string, error) {
	return "legacy", nil
}
func (l *localAdapter) QueryService(ctx context.Context, name string) (*ServiceState, error) {
	return &ServiceState{}, nil
}
func (l *localAdapter) StartService(ctx context.Context, name string) error {
	return nil
}
func (l *localAdapter) StopService(ctx context.Context, name string) error {
	return nil
}
func (l *localAdapter) RestartService(ctx context.Context, name string) error {
	return nil
}

func NewLocalRuntime(logger ...interface{}) Runtime {
	return &localAdapter{}
}
