package kubernetes

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	"github.com/Ajayvtl/devserver/internal/executor/local"
)

type kubernetesAdapter struct {
	id        string
	envID     common.EnvironmentID
	namespace string
	localExec contracts.Executor
}

func New(id string, envID common.EnvironmentID, namespace string) contracts.Executor {
	return &kubernetesAdapter{
		id:        id,
		envID:     envID,
		namespace: namespace,
		localExec: local.New(id+"-local-delegate", envID),
	}
}

func (a *kubernetesAdapter) Execute(ctx context.Context, act *action.Action) (*contracts.ExecutionResult, error) {
	if err := a.Validate(ctx, act); err != nil {
		return nil, err
	}

	podName := act.Arguments[0]

	kubectlTarget, _ := valueobjects.NewReference("kubectl")

	args := []string{"exec", podName, "-n", a.namespace, "--", act.Target.String()}
	if len(act.Arguments) > 1 {
		args = append(args, act.Arguments[1:]...)
	}

	k8sAct := &action.Action{
		Target:    kubectlTarget,
		Arguments: args,
	}

	return a.localExec.Execute(ctx, k8sAct)
}

func (a *kubernetesAdapter) Detect(ctx context.Context) (bool, []contracts.Diagnostic, error) {
	kubectlTarget, _ := valueobjects.NewReference("kubectl")
	k8sAct := &action.Action{
		Target:    kubectlTarget,
		Arguments: []string{"version", "--client"},
	}
	res, err := a.localExec.Execute(ctx, k8sAct)
	if err != nil || res.ExitCode != 0 {
		msg := "kubectl CLI unavailable"
		if err != nil {
			msg = err.Error()
		}
		return false, []contracts.Diagnostic{{Level: "error", Message: msg}}, nil
	}
	return true, nil, nil
}

func (a *kubernetesAdapter) Capabilities() contracts.Capabilities {
	return contracts.Capabilities{
		Kubernetes: &contracts.KubernetesCapability{Supported: true, Version: "1.20+"},
	}
}

func (a *kubernetesAdapter) Validate(ctx context.Context, act *action.Action) error {
	if act == nil || len(act.Arguments) == 0 {
		return contracts.ErrValidationFailed
	}
	return nil
}

func (a *kubernetesAdapter) Health(ctx context.Context) (common.HealthState, error) {
	if ok, _, _ := a.Detect(ctx); ok {
		return common.HealthStateHealthy, nil
	}
	return common.HealthStateUnhealthy, nil
}

func (a *kubernetesAdapter) Metadata() contracts.ExecutorMetadata {
	return contracts.ExecutorMetadata{ID: a.id, EnvironmentID: a.envID, Type: common.ExecutorTypeKubernetes, Version: "1.0"}
}
