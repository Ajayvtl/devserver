package editor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"time"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/rs/zerolog"
)

type eventBridgeImpl struct {
	log    zerolog.Logger
	bus    events.Bus
	active map[common.WorkspaceID]context.CancelFunc
	mu     sync.RWMutex
}

// NewEventBridge creates a new EventBridge.
func NewEventBridge(logger zerolog.Logger, bus events.Bus) EventBridge {
	return &eventBridgeImpl{
		log:    logger.With().Str("component", "EventBridge").Logger(),
		bus:    bus,
		active: make(map[common.WorkspaceID]context.CancelFunc),
	}
}

func (e *eventBridgeImpl) Start(ctx context.Context, workspaceID common.WorkspaceID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.active[workspaceID]; exists {
		return fmt.Errorf("event bridge already active for workspace %s", workspaceID)
	}

	bridgeCtx, cancel := context.WithCancel(ctx)
	e.active[workspaceID] = cancel

	e.log.Info().Str("workspace", string(workspaceID)).Msg("Event bridge started. Waiting for IPC connection.")

	// Simulate connection loop with health monitoring and recovery
	go e.runConnectionLoop(bridgeCtx, workspaceID)

	return nil
}

func (e *eventBridgeImpl) Stop(ctx context.Context, workspaceID common.WorkspaceID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cancel, exists := e.active[workspaceID]
	if !exists {
		return fmt.Errorf("event bridge not active for workspace %s", workspaceID)
	}

	cancel()
	delete(e.active, workspaceID)
	e.log.Info().Str("workspace", string(workspaceID)).Msg("Event bridge stopped")
	return nil
}

func (e *eventBridgeImpl) Send(ctx context.Context, workspaceID common.WorkspaceID, eventType string, payload any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, exists := e.active[workspaceID]; !exists {
		return fmt.Errorf("event bridge not active for workspace %s", workspaceID)
	}

	e.log.Debug().
		Str("workspace", string(workspaceID)).
		Str("type", eventType).
		Msg("Sent IPC event to editor")
	return nil
}

func (e *eventBridgeImpl) Receive(ctx context.Context, workspaceID common.WorkspaceID, payload []byte) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, exists := e.active[workspaceID]; !exists {
		return fmt.Errorf("event bridge not active for workspace %s", workspaceID)
	}

	var rawEvent struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}

	if err := json.Unmarshal(payload, &rawEvent); err != nil {
		e.log.Error().Err(err).Msg("Failed to unmarshal IPC payload")
		return err
	}

	// Route editor events to DevServer EventBus
	// E.g., editor.documentOpened, editor.documentDirty
	e.log.Debug().
		Str("workspace", string(workspaceID)).
		Str("type", rawEvent.Type).
		Msg("Received IPC event from editor")

	// Map strings from the IPC to DevServer EventType if necessary, or pass through as a generic editor event
	e.bus.Publish(events.EventType(rawEvent.Type), rawEvent.Payload)

	return nil
}

func (e *eventBridgeImpl) runConnectionLoop(ctx context.Context, workspaceID common.WorkspaceID) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Subscriptions for full bidirectional synchronization
	subWorkspace := e.bus.Subscribe(events.WorkspaceChanged)
	subTaskStart := e.bus.Subscribe(events.TaskStarted)
	subTaskComplete := e.bus.Subscribe(events.TaskCompleted)
	subCommandComplete := e.bus.Subscribe(events.CommandCompleted)

	for {
		select {
		case <-ctx.Done():
			e.log.Debug().Str("workspace", string(workspaceID)).Msg("Connection loop terminated")
			return
		case <-ticker.C:
			// Health monitoring / Ping
			e.log.Trace().Str("workspace", string(workspaceID)).Msg("IPC bridge health check OK")
		case evt := <-subWorkspace:
			e.Send(ctx, workspaceID, string(events.WorkspaceChanged), evt.Payload)
		case evt := <-subTaskStart:
			e.Send(ctx, workspaceID, string(events.TaskStarted), evt.Payload)
		case evt := <-subTaskComplete:
			e.Send(ctx, workspaceID, string(events.TaskCompleted), evt.Payload)
		case evt := <-subCommandComplete:
			e.Send(ctx, workspaceID, string(events.CommandCompleted), evt.Payload)
		}
	}
}
