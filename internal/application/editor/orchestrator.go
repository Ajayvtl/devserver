package editor

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/events"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/rs/zerolog"

	"github.com/Ajayvtl/devserver/internal/application/envcontext"
	"github.com/Ajayvtl/devserver/internal/state"
)

// Orchestrator coordinates all sub-components of the Editor subsystem.
// It conforms to the runtime.Component interface for unified DevServer lifecycle management.
type Orchestrator struct {
	log              zerolog.Logger
	sessionManager   SessionManager
	processManager   ProcessManager
	extensionManager ExtensionManager
	proxyManager     ProxyManager
	eventBridge      EventBridge
	db               *state.StoreDB

	status rt.Status
	mu     sync.RWMutex
}

// OrchestratorDeps holds dependencies required to construct the Orchestrator.
type OrchestratorDeps struct {
	Logger           zerolog.Logger
	Bus              events.Bus
	SessionManager   SessionManager
	ProcessManager   ProcessManager
	ExtensionManager ExtensionManager
	ProxyManager     ProxyManager
	EventBridge      EventBridge
	StoreDB          *state.StoreDB
}

// NewOrchestrator creates a new instance of the Editor Orchestrator.
func NewOrchestrator(deps OrchestratorDeps) *Orchestrator {
	return &Orchestrator{
		log:              deps.Logger,
		sessionManager:   deps.SessionManager,
		processManager:   deps.ProcessManager,
		extensionManager: deps.ExtensionManager,
		proxyManager:     deps.ProxyManager,
		eventBridge:      deps.EventBridge,
		db:               deps.StoreDB,
		status:           rt.StatusStopped,
	}
}

// Name returns the unique identifier for this runtime component.
func (o *Orchestrator) Name() string {
	return "editor.Orchestrator"
}

// Initialize prepares the component.
func (o *Orchestrator) Initialize(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.status != rt.StatusStopped {
		return fmt.Errorf("editor orchestrator already initialized or running")
	}

	o.status = rt.StatusStarting
	o.log.Info().Msg("Editor orchestrator initialized")
	return nil
}

// Start begins the component's operation.
func (o *Orchestrator) Start(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.status = rt.StatusRunning
	o.log.Info().Msg("Editor orchestrator started")
	return nil
}

// Stop gracefully halts the component and running editor processes.
func (o *Orchestrator) Stop(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.status = rt.StatusStopped
	o.log.Info().Msg("Editor orchestrator stopped")
	return nil
}

// Status returns the current operational state.
func (o *Orchestrator) Status() rt.Status {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.status
}

// Health returns the current health condition.
func (o *Orchestrator) Health() rt.Health {
	return rt.HealthHealthy
}

// ProvisionEditor spins up an editor instance for a workspace, setting up proxy routes and bridges.
func (o *Orchestrator) ProvisionEditor(ctx context.Context, workspaceID common.WorkspaceID) error {
	o.log.Info().Str("workspace", string(workspaceID)).Msg("Provisioning editor instance")

	// Resolve the active environment
	switcher := envcontext.NewSwitcher(o.db)
	envID, _ := switcher.GetActiveEnvironment(ctx, workspaceID)
	ctx = switcher.WithActiveEnvironment(ctx, envID)
	o.log.Debug().Str("environment", string(envID)).Msg("Using active environment for editor session")

	if o.processManager != nil {
		if err := o.processManager.Start(ctx, workspaceID); err != nil {
			return fmt.Errorf("failed to start editor process: %w", err)
		}
	}

	if o.proxyManager != nil {
		targetURL := fmt.Sprintf("http://localhost:8080/workspace/%s", string(workspaceID)) // mock target
		if err := o.proxyManager.RegisterRoute(workspaceID, targetURL); err != nil {
			return fmt.Errorf("failed to register proxy route: %w", err)
		}
	}

	if o.eventBridge != nil {
		if err := o.eventBridge.Start(ctx, workspaceID); err != nil {
			return fmt.Errorf("failed to start event bridge: %w", err)
		}
	}

	return nil
}

// TeardownEditor halts an editor instance and cleans up related routes and bridges.
func (o *Orchestrator) TeardownEditor(ctx context.Context, workspaceID common.WorkspaceID) error {
	o.log.Info().Str("workspace", string(workspaceID)).Msg("Tearing down editor instance")

	var errs []error

	if o.eventBridge != nil {
		if err := o.eventBridge.Stop(ctx, workspaceID); err != nil {
			errs = append(errs, err)
		}
	}

	if o.proxyManager != nil {
		if err := o.proxyManager.RemoveRoute(workspaceID); err != nil {
			errs = append(errs, err)
		}
	}

	if o.processManager != nil {
		if err := o.processManager.Stop(ctx, workspaceID); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during teardown: %v", errs)
	}

	return nil
}
