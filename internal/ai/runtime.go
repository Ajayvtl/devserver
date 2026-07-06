package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/providers"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/rs/zerolog"
)

// Runtime orchestrates LLM dispatch and context gathering for the DevServer AI features.
type Runtime struct {
	log        zerolog.Logger
	workspace  *core.WorkspaceProvider
	indexer    *core.Indexer
	assembler  *ContextAssembler
	dispatcher *Dispatcher

	mu     sync.RWMutex
	status rt.Status
}

// NewRuntime creates a new AI Runtime component.
func NewRuntime(logger zerolog.Logger, workspace *core.WorkspaceProvider, indexer *core.Indexer, pm *providers.Manager) *Runtime {
	return &Runtime{
		log:        logger.With().Str("component", "AIRuntime").Logger(),
		workspace:  workspace,
		indexer:    indexer,
		assembler:  NewContextAssembler(logger, workspace),
		dispatcher: NewDispatcher(logger, pm),
		status:     rt.StatusStopped,
	}
}

// Name returns the component's name.
func (r *Runtime) Name() string {
	return "ai.Runtime"
}

// Initialize sets up the AI Runtime.
func (r *Runtime) Initialize(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = rt.StatusStarting
	r.log.Debug().Msg("Initializing AI Runtime")
	return nil
}

// Start begins processing context streams and dispatching LLM queries.
func (r *Runtime) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.workspace == nil || r.indexer == nil {
		r.status = rt.StatusError
		return fmt.Errorf("ai runtime requires workspace provider and indexer")
	}

	r.status = rt.StatusRunning
	r.log.Info().Msg("AI Runtime started")
	return nil
}

// Stop halts all active context tracking and LLM dispatches.
func (r *Runtime) Stop(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status = rt.StatusStopped
	r.log.Info().Msg("AI Runtime stopped")
	return nil
}

// Status returns the current operational state.
func (r *Runtime) Status() rt.Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.status
}

// Health checks if the AI Runtime and its dependencies are functional.
func (r *Runtime) Health() rt.Health {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.status != rt.StatusRunning {
		return rt.HealthUnhealthy
	}
	return rt.HealthHealthy
}

// Infer processes an inference request, orchestrating context gathering and dispatching it to the LLM.
func (r *Runtime) Infer(ctx context.Context, workspaceID common.WorkspaceID, req InferenceRequest) (*InferenceResponse, error) {
	if r.Status() != rt.StatusRunning {
		return nil, fmt.Errorf("ai runtime is not running")
	}

	// 1. Gather Context
	wsContext, err := r.assembler.Assemble(ctx, workspaceID)
	if err != nil {
		r.log.Warn().Err(err).Msg("Failed to assemble full workspace context, continuing with partial context")
	}

	// 2. Attach Context
	req.Context = wsContext

	// 3. Dispatch Request
	return r.dispatcher.Dispatch(ctx, req)
}
