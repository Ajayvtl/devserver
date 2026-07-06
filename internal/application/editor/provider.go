package editor

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/providers"
)

// VSCodeExtensionProvider wraps a VS Code extension as a DevServer Provider.
type VSCodeExtensionProvider struct {
	extensionID string
	metadata    providers.ProviderMetadata
	extManager  ExtensionManager
}

// NewVSCodeExtensionProvider creates a new provider for a VS Code extension.
func NewVSCodeExtensionProvider(extensionID, name, description, version string, extManager ExtensionManager) providers.Provider {
	return &VSCodeExtensionProvider{
		extensionID: extensionID,
		metadata: providers.ProviderMetadata{
			Name:        name,
			Description: description,
			Version:     version,
			Tags:        []string{"editor-extension", "vscode"},
		},
		extManager: extManager,
	}
}

func (p *VSCodeExtensionProvider) Metadata() providers.ProviderMetadata {
	return p.metadata
}

func (p *VSCodeExtensionProvider) Detect(ctx context.Context) (bool, error) {
	list, err := p.extManager.List(ctx)
	if err != nil {
		return false, err
	}
	for _, ext := range list {
		if ext == p.extensionID {
			return true, nil
		}
	}
	return false, nil
}

func (p *VSCodeExtensionProvider) Version(ctx context.Context) (string, error) {
	installed, err := p.Detect(ctx)
	if err != nil {
		return "", err
	}
	if !installed {
		return "", fmt.Errorf("extension not installed")
	}
	return p.metadata.Version, nil
}

func (p *VSCodeExtensionProvider) Health(ctx context.Context) error {
	installed, err := p.Detect(ctx)
	if err != nil {
		return err
	}
	if !installed {
		return fmt.Errorf("extension not installed")
	}
	return nil
}

func (p *VSCodeExtensionProvider) Status(ctx context.Context) (providers.ProviderStatus, error) {
	installed, err := p.Detect(ctx)
	if err != nil {
		return providers.StatusFailed, err
	}
	if installed {
		return providers.StatusRunning, nil // extensions are "running" if the editor is running
	}
	return providers.StatusNotInstalled, nil
}

func (p *VSCodeExtensionProvider) Info(ctx context.Context) (providers.ProviderInfo, error) {
	status, err := p.Status(ctx)
	installed := status == providers.StatusRunning

	return providers.ProviderInfo{
		Name:      p.metadata.Name,
		Version:   p.metadata.Version,
		Installed: installed,
		State: providers.ProviderState{
			Status: status,
			Health: providers.HealthHealthy, // Simplification
		},
		Metrics:      providers.ProviderMetrics{},
		Capabilities: p.Capabilities(),
	}, err
}

func (p *VSCodeExtensionProvider) Capabilities() providers.ProviderCapabilities {
	return providers.ProviderCapabilities{
		Detect:    true,
		Version:   true,
		Health:    true,
		Status:    true,
		Install:   true,
		Uninstall: true,
	}
}

func (p *VSCodeExtensionProvider) Install(ctx context.Context) error {
	return p.extManager.Install(ctx, p.extensionID)
}

func (p *VSCodeExtensionProvider) Update(ctx context.Context) error {
	// Re-installing typically updates it
	return p.extManager.Install(ctx, p.extensionID)
}

func (p *VSCodeExtensionProvider) Uninstall(ctx context.Context) error {
	return p.extManager.Remove(ctx, p.extensionID)
}

func (p *VSCodeExtensionProvider) Start(ctx context.Context) error {
	return fmt.Errorf("extensions are managed by the editor lifecycle")
}

func (p *VSCodeExtensionProvider) Stop(ctx context.Context) error {
	return fmt.Errorf("extensions are managed by the editor lifecycle")
}

func (p *VSCodeExtensionProvider) Restart(ctx context.Context) error {
	return fmt.Errorf("extensions are managed by the editor lifecycle")
}

func (p *VSCodeExtensionProvider) Configure(ctx context.Context, config map[string]any) error {
	return fmt.Errorf("configuration not yet supported for extensions")
}

func (p *VSCodeExtensionProvider) Validate(ctx context.Context) error {
	return nil
}
