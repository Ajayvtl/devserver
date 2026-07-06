package editor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

type processManagerImpl struct {
	log    zerolog.Logger
	engine tasks.Engine
	status map[common.WorkspaceID]rt.Status
	mu     sync.RWMutex
}

// NewProcessManager creates a new ProcessManager backed by the Task Engine.
func NewProcessManager(logger zerolog.Logger, engine tasks.Engine) ProcessManager {
	return &processManagerImpl{
		log:    logger.With().Str("component", "ProcessManager").Logger(),
		engine: engine,
		status: make(map[common.WorkspaceID]rt.Status),
	}
}

func (p *processManagerImpl) Start(ctx context.Context, workspaceID common.WorkspaceID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status[workspaceID] == rt.StatusRunning {
		return fmt.Errorf("editor process already running for workspace %s", workspaceID)
	}

	p.status[workspaceID] = rt.StatusStarting

	task := &tasks.Task{
		ID:          "editor-startup-" + string(workspaceID),
		WorkspaceID: string(workspaceID),
		Name:        "Start OpenVSCode Server",
		Type:        "editor.process.start",
		Priority:    tasks.TaskPriorityHigh,
		Timeout:     time.Minute * 5,
		Payload: map[string]any{
			"port": 8080,
			"host": "127.0.0.1",
		},
	}

	// Start async
	if _, err := p.engine.Submit(task); err != nil {
		p.status[workspaceID] = rt.StatusError
		return fmt.Errorf("failed to submit editor start task: %w", err)
	}

	p.status[workspaceID] = rt.StatusRunning
	p.log.Info().Str("workspace", string(workspaceID)).Msg("Editor process started")
	return nil
}

func (p *processManagerImpl) Stop(ctx context.Context, workspaceID common.WorkspaceID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status[workspaceID] != rt.StatusRunning {
		return nil // Nothing to stop
	}

	task := &tasks.Task{
		ID:          "editor-shutdown-" + string(workspaceID),
		WorkspaceID: string(workspaceID),
		Name:        "Stop OpenVSCode Server",
		Type:        "editor.process.stop",
		Priority:    tasks.TaskPriorityHigh,
		Timeout:     time.Second * 30,
	}

	if _, err := p.engine.Submit(task); err != nil {
		return fmt.Errorf("failed to submit editor stop task: %w", err)
	}

	p.status[workspaceID] = rt.StatusStopped
	p.log.Info().Str("workspace", string(workspaceID)).Msg("Editor process stopped")
	return nil
}

func (p *processManagerImpl) Status(ctx context.Context, workspaceID common.WorkspaceID) rt.Status {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if st, ok := p.status[workspaceID]; ok {
		return st
	}
	return rt.StatusStopped
}
