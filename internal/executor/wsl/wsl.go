package wsl

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	"github.com/Ajayvtl/devserver/internal/executor/local"
)

type wslAdapter struct {
	id         string
	distro     string
	localExec  contracts.Executor
}

func New(id string, distro string) contracts.Executor {
	return &wslAdapter{
		id:        id,
		distro:    distro,
		localExec: local.New(id + "-local-delegate"),
	}
}

func (a *wslAdapter) Execute(ctx context.Context, act *action.Action) (*contracts.ExecutionResult, error) {
	if err := a.Validate(ctx, act); err != nil {
		return nil, err
	}

	wslTarget, _ := valueobjects.NewReference("wsl.exe")
	
	args := []string{"-d", a.distro, "-e", act.Target.String()}
	args = append(args, act.Arguments...)

	wslAct := &action.Action{
		Target:    wslTarget,
		Arguments: args,
	}

	return a.localExec.Execute(ctx, wslAct)
}

func (a *wslAdapter) Detect(ctx context.Context) (bool, []contracts.Diagnostic, error) {
	// Simple detect by running wsl -l
	wslTarget, _ := valueobjects.NewReference("wsl.exe")
	wslAct := &action.Action{
		Target:    wslTarget,
		Arguments: []string{"-l"},
	}
	res, err := a.localExec.Execute(ctx, wslAct)
	if err != nil || res.ExitCode != 0 {
		msg := "wsl CLI unavailable"
		if err != nil {
			msg = err.Error()
		}
		return false, []contracts.Diagnostic{{Level: "error", Message: msg}}, nil
	}
	return true, nil, nil
}

func (a *wslAdapter) Capabilities() contracts.Capabilities {
	return contracts.Capabilities{
		Shell: &contracts.ShellCapability{Supported: true, Type: "bash"},
	}
}

func (a *wslAdapter) Validate(ctx context.Context, act *action.Action) error {
	if act == nil {
		return contracts.ErrValidationFailed
	}
	return nil
}

func (a *wslAdapter) Health(ctx context.Context) (common.HealthState, error) {
	if ok, _, _ := a.Detect(ctx); ok {
		return common.HealthStateHealthy, nil
	}
	return common.HealthStateUnhealthy, nil
}

func (a *wslAdapter) Metadata() contracts.ExecutorMetadata {
	return contracts.ExecutorMetadata{ID: a.id, Type: common.ExecutorTypeWSL, Version: "1.0"}
}
