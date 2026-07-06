package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/executor/legacy"
)

// RedisProvider manages a Redis service using the legacy.Runtime.
type RedisProvider struct {
	Runtime legacy.Runtime
	Meta    ProviderMetadata
}

func NewRedisProvider(runtime legacy.Runtime) *RedisProvider {
	return &RedisProvider{
		Runtime: runtime,
		Meta: ProviderMetadata{
			Name:        "redis",
			Description: "In-memory data structure store",
			Version:     "latest",
		},
	}
}

func (p *RedisProvider) Metadata() ProviderMetadata {
	return p.Meta
}

func (p *RedisProvider) Detect(ctx context.Context) (bool, error) {
	return p.Runtime.DetectBinary(ctx, "redis-cli")
}

func (p *RedisProvider) Version(ctx context.Context) (string, error) {
	installed, err := p.Detect(ctx)
	if !installed || err != nil {
		return "unknown", nil
	}
	out, err := p.Runtime.QueryVersion(ctx, "redis-cli", "--version")
	if err != nil {
		return "unknown", fmt.Errorf("failed to get redis version: %w", err)
	}
	return out, nil
}

func (p *RedisProvider) Health(ctx context.Context) error {
	installed, _ := p.Detect(ctx)
	if !installed {
		return fmt.Errorf("redis is not installed")
	}
	
	// Check if we can ping it.
	res, err := p.Runtime.Execute(ctx, "redis-cli", "ping")
	if err != nil || res.ExitCode != 0 {
		return fmt.Errorf("redis health check failed: %v", res.Error)
	}
	if !strings.Contains(res.Stdout, "PONG") {
		return fmt.Errorf("redis returned unexpected ping response: %s", res.Stdout)
	}
	return nil
}

func (p *RedisProvider) Status(ctx context.Context) (ProviderStatus, error) {
	state, err := p.Runtime.QueryService(ctx, "redis")
	if err != nil {
		if err == legacy.ErrServiceNotInstalled {
			return StatusNotInstalled, nil
		}
		return StatusFailed, err
	}
	if state.Running {
		return StatusRunning, nil
	}
	return StatusStopped, nil
}

func (p *RedisProvider) Info(ctx context.Context) (ProviderInfo, error) {
	state, err := p.Runtime.QueryService(ctx, "redis")
	if err != nil && err != legacy.ErrServiceNotInstalled {
		state = &legacy.ServiceState{Status: "unknown"}
	} else if state == nil {
		state = &legacy.ServiceState{Status: "not_installed"}
	}

	status, _ := p.Status(ctx)
	health := HealthUnknown
	if status == StatusRunning {
		if err := p.Health(ctx); err == nil {
			health = HealthHealthy
		} else {
			health = HealthDegraded
		}
	} else if status == StatusStopped {
		health = HealthNotApplicable
	} else if status == StatusNotInstalled {
		health = HealthNotApplicable
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
		Metrics: ProviderMetrics{
			PID:    state.PID,
			Uptime: state.Uptime.String(),
		},
		Capabilities: p.Capabilities(),
	}, nil
}

func (p *RedisProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Detect:  true,
		Version: true,
		Health:  true,
		Status:  true,
		Start:   true,
		Stop:    true,
		Restart: true,
	}
}

func (p *RedisProvider) Install(ctx context.Context) error {
	return fmt.Errorf("install not implemented for redis")
}

func (p *RedisProvider) Update(ctx context.Context) error {
	return fmt.Errorf("update not implemented for redis")
}

func (p *RedisProvider) Uninstall(ctx context.Context) error {
	return fmt.Errorf("uninstall not implemented for redis")
}

func (p *RedisProvider) Start(ctx context.Context) error {
	return p.Runtime.StartService(ctx, "redis")
}

func (p *RedisProvider) Stop(ctx context.Context) error {
	return p.Runtime.StopService(ctx, "redis")
}

func (p *RedisProvider) Restart(ctx context.Context) error {
	return p.Runtime.RestartService(ctx, "redis")
}

func (p *RedisProvider) Configure(ctx context.Context, config map[string]any) error {
	return fmt.Errorf("configure not implemented for redis")
}

func (p *RedisProvider) Validate(ctx context.Context) error {
	return p.Health(ctx)
}
