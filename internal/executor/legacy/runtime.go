package legacy

import (
	"context"
	"errors"
	"time"
)

// Deprecated: ErrUnsupportedOperation is maintained for legacy compatibility. Use contracts.ErrNotSupported instead.
var ErrUnsupportedOperation = errors.New("operation not supported")
var ErrServiceNotInstalled = errors.New("service not installed")

type LegacyExecutionResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

type ServiceState struct {
	PID     int
	Status  string
	Uptime  time.Duration
	Running bool
}

// Deprecated: Runtime is maintained for legacy compatibility. Use contracts.Executor instead.
type Runtime interface {
	Runner
	Execute(ctx context.Context, cmd string, args ...string) (*LegacyExecutionResult, error)
	StartProcess(ctx context.Context, cmd string, args ...string) error
	StopProcess(ctx context.Context, name string) error
	CheckHealth(ctx context.Context, name string) bool
	DetectBinary(ctx context.Context, name string) (bool, error)
	QueryVersion(ctx context.Context, cmd string, args ...string) (string, error)
	QueryService(ctx context.Context, name string) (*ServiceState, error)
	StartService(ctx context.Context, name string) error
	StopService(ctx context.Context, name string) error
	RestartService(ctx context.Context, name string) error
}

// Deprecated: Runner is maintained for legacy compatibility. Use contracts.Executor instead.
type Runner interface {
	Run(task interface{}) error
}
