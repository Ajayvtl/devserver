package editor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

// extensionManagerImpl implements the ExtensionManager interface.
type extensionManagerImpl struct {
	log        zerolog.Logger
	engine     tasks.Engine
	extensions map[string]bool // state tracking
	mu         sync.RWMutex
}

// NewExtensionManager creates a new instance of ExtensionManager.
func NewExtensionManager(logger zerolog.Logger, engine tasks.Engine) ExtensionManager {
	return &extensionManagerImpl{
		log:        logger.With().Str("component", "ExtensionManager").Logger(),
		engine:     engine,
		extensions: make(map[string]bool),
	}
}

// Install simulates the installation of a VS Code extension.
func (e *extensionManagerImpl) Install(ctx context.Context, extensionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.extensions[extensionID] {
		return fmt.Errorf("extension %s is already installed", extensionID)
	}

	task := &tasks.Task{
		ID:       "install-ext-" + extensionID,
		Name:     "Install VS Code Extension",
		Type:     "editor.extension.install",
		Priority: tasks.TaskPriorityHigh,
		Timeout:  time.Minute * 2,
		Payload: map[string]any{
			"extensionId": extensionID,
		},
	}

	if _, err := e.engine.Submit(task); err != nil {
		return fmt.Errorf("failed to submit extension install task: %w", err)
	}

	e.extensions[extensionID] = true
	e.log.Info().Str("extension_id", extensionID).Msg("Extension installed successfully")
	return nil
}

// Remove simulates the removal of a VS Code extension.
func (e *extensionManagerImpl) Remove(ctx context.Context, extensionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.extensions[extensionID] {
		return fmt.Errorf("extension %s is not installed", extensionID)
	}

	task := &tasks.Task{
		ID:       "uninstall-ext-" + extensionID,
		Name:     "Uninstall VS Code Extension",
		Type:     "editor.extension.uninstall",
		Priority: tasks.TaskPriorityHigh,
		Timeout:  time.Minute * 2,
		Payload: map[string]any{
			"extensionId": extensionID,
		},
	}

	if _, err := e.engine.Submit(task); err != nil {
		return fmt.Errorf("failed to submit extension uninstall task: %w", err)
	}

	delete(e.extensions, extensionID)
	e.log.Info().Str("extension_id", extensionID).Msg("Extension removed successfully")
	return nil
}

// List returns all currently installed extensions.
func (e *extensionManagerImpl) List(ctx context.Context) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var list []string
	for ext := range e.extensions {
		list = append(list, ext)
	}
	return list, nil
}
