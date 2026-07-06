package providers

import (
	"context"
)

// ProviderStatus represents the runtime state of a provider.
type ProviderStatus string

const (
	StatusRunning      ProviderStatus = "running"
	StatusStopped      ProviderStatus = "stopped"
	StatusInstalling   ProviderStatus = "installing"
	StatusUpdating     ProviderStatus = "updating"
	StatusNotInstalled ProviderStatus = "not_installed"
	StatusDisabled     ProviderStatus = "disabled"
	StatusFailed       ProviderStatus = "failed"
)

// ProviderHealth represents the health condition of a provider.
type ProviderHealth string

const (
	HealthHealthy       ProviderHealth = "healthy"
	HealthDegraded      ProviderHealth = "degraded"
	HealthFailed        ProviderHealth = "failed"
	HealthNotApplicable ProviderHealth = "not_applicable"
	HealthUnknown       ProviderHealth = "unknown"
)

// ProviderState represents the active operational condition.
type ProviderState struct {
	Status ProviderStatus `json:"status"`
	Health ProviderHealth `json:"health"`
}

// ProviderMetrics represents runtime metrics.
type ProviderMetrics struct {
	PID    int    `json:"pid,omitempty"`
	Uptime string `json:"uptime,omitempty"`
}

// ProviderInfo represents the canonical state and capabilities of a provider.
type ProviderInfo struct {
	Name         string               `json:"name"`
	Version      string               `json:"version,omitempty"`
	Installed    bool                 `json:"installed"`
	State        ProviderState        `json:"state"`
	Metrics      ProviderMetrics      `json:"metrics"`
	Capabilities ProviderCapabilities `json:"capabilities"`
}

// ProviderMetadata describes the provider.
type ProviderMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

// Provider represents a managed dependency (Node, Redis, Docker, etc.).
type Provider interface {
	// Metadata returns information about the provider.
	Metadata() ProviderMetadata

	// Detect checks if the provider is installed on the host system.
	Detect(ctx context.Context) (bool, error)

	// Version returns the installed version.
	Version(ctx context.Context) (string, error)

	// Health checks if the provider is functioning correctly.
	Health(ctx context.Context) error

	// Status returns the current runtime status (e.g. running, stopped).
	Status(ctx context.Context) (ProviderStatus, error)

	// Info returns detailed state of a background service.
	Info(ctx context.Context) (ProviderInfo, error)

	// Capabilities returns the supported operations.
	Capabilities() ProviderCapabilities

	// Lifecycle management
	Install(ctx context.Context) error
	Update(ctx context.Context) error
	Uninstall(ctx context.Context) error

	// Service management
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error

	// Configuration
	Configure(ctx context.Context, config map[string]any) error
	Validate(ctx context.Context) error
}

// ProviderCapabilities defines which actions a provider supports.
type ProviderCapabilities struct {
	Detect        bool `json:"detect"`
	Version       bool `json:"version"`
	Health        bool `json:"health"`
	Status        bool `json:"status"`
	Start         bool `json:"start"`
	Stop          bool `json:"stop"`
	Restart       bool `json:"restart"`
	Install       bool `json:"install"`
	Update        bool `json:"update"`
	Uninstall     bool `json:"uninstall"`
	Configure     bool `json:"configure"`
	Logs          bool `json:"logs"`
	Metrics       bool `json:"metrics"`
	Shell         bool `json:"shell"`
	OpenUI        bool `json:"open_ui"`
	Documentation bool `json:"documentation"`
}


