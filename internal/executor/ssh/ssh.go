package ssh

import (
	"bytes"
	"context"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	internalUtils "github.com/Ajayvtl/devserver/internal/executor/internal"
	"golang.org/x/crypto/ssh"
)

type sshAdapter struct {
	id     string
	envID  common.EnvironmentID
	host   string
	config *ssh.ClientConfig
	client *ssh.Client
}

func New(id string, envID common.EnvironmentID, host string, config *ssh.ClientConfig) contracts.Executor {
	return &sshAdapter{
		id:     id,
		envID:  envID,
		host:   host,
		config: config,
	}
}

func (a *sshAdapter) Execute(ctx context.Context, act *action.Action) (*contracts.ExecutionResult, error) {
	if err := a.Validate(ctx, act); err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", a.host, a.config)
	if err != nil {
		return nil, internalUtils.NormalizeError(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, internalUtils.NormalizeError(err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	cmdName := act.Target.String()
	for _, arg := range act.Arguments {
		cmdName += " " + arg // simple naive quoting for tests
	}

	start := time.Now()
	err = session.Run(cmdName)
	duration := time.Since(start)

	res := &contracts.ExecutionResult{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: 0,
		Duration: duration,
	}

	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			res.ExitCode = exitError.ExitStatus()
		} else {
			res.ExitCode = -1
		}
	}

	return res, internalUtils.NormalizeError(err)
}

func (a *sshAdapter) Detect(ctx context.Context) (bool, []contracts.Diagnostic, error) {
	client, err := ssh.Dial("tcp", a.host, a.config)
	if err != nil {
		return false, []contracts.Diagnostic{{Level: "error", Message: err.Error()}}, nil
	}
	client.Close()
	return true, nil, nil
}

func (a *sshAdapter) Capabilities() contracts.Capabilities {
	return contracts.Capabilities{
		Shell: &contracts.ShellCapability{Supported: true, Type: "bash"},
	}
}

func (a *sshAdapter) Validate(ctx context.Context, act *action.Action) error {
	if act == nil {
		return contracts.ErrValidationFailed
	}
	return nil
}

func (a *sshAdapter) Health(ctx context.Context) (common.HealthState, error) {
	if ok, _, _ := a.Detect(ctx); ok {
		return common.HealthStateHealthy, nil
	}
	return common.HealthStateUnhealthy, nil
}

func (a *sshAdapter) Metadata() contracts.ExecutorMetadata {
	return contracts.ExecutorMetadata{ID: a.id, EnvironmentID: a.envID, Type: common.ExecutorTypeSSH, Version: "1.0"}
}
