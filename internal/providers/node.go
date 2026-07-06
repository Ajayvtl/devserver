package providers

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/executor/legacy"
)

// NodeProvider manages the Node.js runtime environment using the legacy.Runtime.
type NodeProvider struct {
	Runtime legacy.Runtime
	Meta    ProviderMetadata
}

func NewNodeProvider(runtime legacy.Runtime) *NodeProvider {
	return &NodeProvider{
		Runtime: runtime,
		Meta: ProviderMetadata{
			Name:        "node",
			Description: "Node.js JavaScript runtime",
			Version:     "latest",
		},
	}
}

func (p *NodeProvider) Metadata() ProviderMetadata {
	return p.Meta
}

func (p *NodeProvider) Detect(ctx context.Context) (bool, error) {
	return p.Runtime.DetectBinary(ctx, "node")
}

func (p *NodeProvider) Version(ctx context.Context) (string, error) {
	installed, err := p.Detect(ctx)
	if !installed || err != nil {
		return "unknown", nil
	}
	out, err := p.Runtime.QueryVersion(ctx, "node", "--version")
	if err != nil {
		return "unknown", fmt.Errorf("failed to get node version: %w", err)
	}
	return out, nil
}

func (p *NodeProvider) Health(ctx context.Context) error {
	installed, _ := p.Detect(ctx)
	if !installed {
		return fmt.Errorf("node is not installed")
	}
	return nil // Node is healthy if it exists on the system
}

func (p *NodeProvider) Status(ctx context.Context) (ProviderStatus, error) {
	installed, _ := p.Detect(ctx)
	if !installed {
		return StatusNotInstalled, nil
	}
	return StatusStopped, nil // Node itself is a runtime, not a service
}

func (p *NodeProvider) Info(ctx context.Context) (ProviderInfo, error) {
	status, _ := p.Status(ctx)
	health := HealthUnknown
	if status == StatusStopped {
		health = HealthNotApplicable
	} else if status == StatusNotInstalled {
		health = HealthNotApplicable
	} else if err := p.Health(ctx); err == nil {
		health = HealthHealthy
	} else {
		health = HealthFailed
	}

	version, _ := p.Version(ctx)
	installed, _ := p.Detect(ctx)

	if !installed {
		version = ""
	}

	return ProviderInfo{
		Name:      p.Meta.Name,
		Version:   version,
		Installed: installed,
		State: ProviderState{
			Status: status,
			Health: health,
		},
		Metrics:      ProviderMetrics{},
		Capabilities: p.Capabilities(),
	}, nil
}

func (p *NodeProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Detect:  true,
		Version: true,
		Health:  true,
		Status:  true,
		Start:   false,
		Stop:    false,
		Restart: false,
	}
}

func (p *NodeProvider) Install(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Update(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Uninstall(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Start(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Stop(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Restart(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Configure(ctx context.Context, config map[string]any) error {
	return legacy.ErrUnsupportedOperation
}

func (p *NodeProvider) Validate(ctx context.Context) error {
	return p.Health(ctx)
}
