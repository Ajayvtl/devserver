package local

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	internalUtils "github.com/Ajayvtl/devserver/internal/executor/internal"
)

type localAdapter struct {
	id      string
	version string
}

func New(id string) contracts.Executor {
	return &localAdapter{
		id:      id,
		version: "1.0.0",
	}
}

func (l *localAdapter) Execute(ctx context.Context, act *action.Action) (*contracts.ExecutionResult, error) {
	if err := l.Validate(ctx, act); err != nil {
		return nil, err
	}

	cmdName := act.Target.String()
	args := act.Arguments

	cmd := exec.CommandContext(ctx, cmdName, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	res := &contracts.ExecutionResult{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		Duration: duration,
		ExitCode: 0,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			res.ExitCode = exitError.ExitCode()
		} else {
			res.ExitCode = -1
		}
	}

	return res, internalUtils.NormalizeError(err)
}

func (l *localAdapter) Detect(ctx context.Context) (bool, []contracts.Diagnostic, error) {
	return true, nil, nil // Local adapter is generally always available
}

func (l *localAdapter) Capabilities() contracts.Capabilities {
	return contracts.Capabilities{
		Shell: &contracts.ShellCapability{Supported: true, Type: "os_default"},
		Filesystem: &contracts.FilesystemCapability{Supported: true},
		Process: &contracts.ProcessCapability{Supported: true},
	}
}

func (l *localAdapter) Validate(ctx context.Context, act *action.Action) error {
	if act == nil {
		return contracts.ErrValidationFailed
	}
	return nil
}

func (l *localAdapter) Health(ctx context.Context) (common.HealthState, error) {
	if ok, _, _ := l.Detect(ctx); ok {
		return common.HealthStateHealthy, nil
	}
	return common.HealthStateUnhealthy, nil
}

func (l *localAdapter) Metadata() contracts.ExecutorMetadata {
	return contracts.ExecutorMetadata{
		ID:      l.id,
		Type:    common.ExecutorTypeLocal,
		Version: l.version,
	}
}
