package providers

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/executor/legacy"
)

// LocalProvider implements Provider for locally installed services.
type LocalProvider struct {
	Runtime     legacy.Runtime
	Meta        ProviderMetadata
	CheckCmd    string
	CheckArgs   []string
	StartCmd    string
	StartArgs   []string
	StopCmd     string
	StopArgs    []string
	RestartCmd  string
	RestartArgs []string
}

func (p *LocalProvider) Metadata() ProviderMetadata {
	return p.Meta
}

func (p *LocalProvider) Detect(ctx context.Context) (bool, error) {
	if p.CheckCmd == "" {
		return false, nil
	}
	res, err := p.Runtime.Execute(ctx, "which", p.CheckCmd)
	if err != nil || res.ExitCode != 0 {
		return false, nil
	}
	return true, nil
}

func (p *LocalProvider) Version(ctx context.Context) (string, error) {
	if p.CheckCmd == "" {
		return "unknown", nil
	}
	res, err := p.Runtime.Execute(ctx, p.CheckCmd, p.CheckArgs...)
	if err != nil || res.ExitCode != 0 {
		return "unknown", fmt.Errorf("failed to get version: %v", res.Error)
	}
	return res.Stdout, nil
}

func (p *LocalProvider) Health(ctx context.Context) error {
	detected, _ := p.Detect(ctx)
	if !detected {
		return fmt.Errorf("provider %s is not installed", p.Meta.Name)
	}
	return nil
}

func (p *LocalProvider) Status(ctx context.Context) (ProviderStatus, error) {
	detected, _ := p.Detect(ctx)
	if !detected {
		return StatusNotInstalled, nil
	}
	// Simplified logic for local providers: if it exists, assume it can be managed.
	return StatusStopped, nil
}

func (p *LocalProvider) Info(ctx context.Context) (ProviderInfo, error) {
	status, _ := p.Status(ctx)
	health := HealthUnknown
	if status == StatusNotInstalled {
		health = HealthNotApplicable
	} else if err := p.Health(ctx); err == nil {
		health = HealthHealthy
	}
	version, _ := p.Version(ctx)
	installed, _ := p.Detect(ctx)

	if !installed {
		version = ""
	}

	return ProviderInfo{
		Name:         p.Meta.Name,
		Version:      version,
		Installed:    installed,
		State: ProviderState{
			Status: status,
			Health: health,
		},
		Metrics: ProviderMetrics{},
		Capabilities: p.Capabilities(),
	}, nil
}

func (p *LocalProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Detect:  p.CheckCmd != "",
		Version: p.CheckCmd != "",
		Health:  p.CheckCmd != "",
		Status:  p.CheckCmd != "",
		Start:   p.StartCmd != "",
		Stop:    p.StopCmd != "",
		Restart: p.RestartCmd != "",
	}
}

func (p *LocalProvider) Install(ctx context.Context) error {
	return fmt.Errorf("install not implemented for %s", p.Meta.Name)
}

func (p *LocalProvider) Update(ctx context.Context) error {
	return fmt.Errorf("update not implemented for %s", p.Meta.Name)
}

func (p *LocalProvider) Uninstall(ctx context.Context) error {
	return fmt.Errorf("uninstall not implemented for %s", p.Meta.Name)
}

func (p *LocalProvider) Start(ctx context.Context) error {
	if p.StartCmd == "" {
		return fmt.Errorf("start not supported for %s", p.Meta.Name)
	}
	return p.Runtime.StartProcess(ctx, p.StartCmd, p.StartArgs...)
}

func (p *LocalProvider) Stop(ctx context.Context) error {
	if p.StopCmd == "" {
		return fmt.Errorf("stop not supported for %s", p.Meta.Name)
	}
	res, err := p.Runtime.Execute(ctx, p.StopCmd, p.StopArgs...)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return res.Error
	}
	return nil
}

func (p *LocalProvider) Restart(ctx context.Context) error {
	if p.RestartCmd == "" {
		// Fallback to Stop() then Start()
		if err := p.Stop(ctx); err != nil {
			// ignore stop errors
		}
		return p.Start(ctx)
	}
	res, err := p.Runtime.Execute(ctx, p.RestartCmd, p.RestartArgs...)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return res.Error
	}
	return nil
}

func (p *LocalProvider) Configure(ctx context.Context, config map[string]any) error {
	return fmt.Errorf("configure not implemented for %s", p.Meta.Name)
}

func (p *LocalProvider) Validate(ctx context.Context) error {
	return p.Health(ctx)
}
