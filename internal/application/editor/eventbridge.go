package editor

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

type eventBridgeImpl struct {
	log    zerolog.Logger
	active map[common.WorkspaceID]bool
	mu     sync.RWMutex
}

// NewEventBridge creates a new EventBridge.
func NewEventBridge(logger zerolog.Logger) EventBridge {
	return &eventBridgeImpl{
		log:    logger.With().Str("component", "EventBridge").Logger(),
		active: make(map[common.WorkspaceID]bool),
	}
}

func (e *eventBridgeImpl) Start(ctx context.Context, workspaceID common.WorkspaceID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.active[workspaceID] {
		return fmt.Errorf("event bridge already active for workspace %s", workspaceID)
	}
	e.active[workspaceID] = true
	e.log.Debug().Str("workspace", string(workspaceID)).Msg("Event bridge started")
	return nil
}

func (e *eventBridgeImpl) Stop(ctx context.Context, workspaceID common.WorkspaceID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.active[workspaceID] {
		return fmt.Errorf("event bridge not active for workspace %s", workspaceID)
	}
	delete(e.active, workspaceID)
	e.log.Debug().Str("workspace", string(workspaceID)).Msg("Event bridge stopped")
	return nil
}
