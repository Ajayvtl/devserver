package docker

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	"github.com/Ajayvtl/devserver/internal/executor/local"
)

type dockerAdapter struct {
	id        string
	envID     common.EnvironmentID
	localExec contracts.Executor
}

func New(id string, envID common.EnvironmentID) contracts.Executor {
	return &dockerAdapter{
		id:        id,
		envID:     envID,
		localExec: local.New(id+"-local-delegate", envID),
	}
}

func (a *dockerAdapter) Execute(ctx context.Context, act *action.Action) (*contracts.ExecutionResult, error) {
	if err := a.Validate(ctx, act); err != nil {
		return nil, err
	}

	containerID := "unknown"
	if len(act.Arguments) > 0 {
		containerID = act.Arguments[0]
	}

	dockerTarget, _ := valueobjects.NewReference("docker")

	args := []string{"exec", containerID, "-i", "--", act.Target.String()}
	if len(act.Arguments) > 1 {
		args = append(args, act.Arguments[1:]...)
	}

	dockerAct := &action.Action{
		Target:    dockerTarget,
		Arguments: args,
	}

	return a.localExec.Execute(ctx, dockerAct)
}

func (a *dockerAdapter) Detect(ctx context.Context) (bool, []contracts.Diagnostic, error) {
	dockerTarget, _ := valueobjects.NewReference("docker")
	dockerAct := &action.Action{
		Target:    dockerTarget,
		Arguments: []string{"info"},
	}
	res, err := a.localExec.Execute(ctx, dockerAct)
	if err != nil || res.ExitCode != 0 {
		msg := "docker CLI unavailable"
		if err != nil {
			msg = err.Error()
		}
		return false, []contracts.Diagnostic{{Level: "error", Message: msg}}, nil
	}
	return true, nil, nil
}

func (a *dockerAdapter) Capabilities() contracts.Capabilities {
	return contracts.Capabilities{
		Docker: &contracts.DockerCapability{Supported: true, Version: "20.10+"},
	}
}

func (a *dockerAdapter) Validate(ctx context.Context, act *action.Action) error {
	if act == nil || len(act.Arguments) == 0 {
		return contracts.ErrValidationFailed
	}
	return nil
}

func (a *dockerAdapter) Health(ctx context.Context) (common.HealthState, error) {
	if ok, _, _ := a.Detect(ctx); ok {
		return common.HealthStateHealthy, nil
	}
	return common.HealthStateUnhealthy, nil
}

func (a *dockerAdapter) Metadata() contracts.ExecutorMetadata {
	return contracts.ExecutorMetadata{ID: a.id, EnvironmentID: a.envID, Type: common.ExecutorTypeDocker, Version: "1.0"}
}
