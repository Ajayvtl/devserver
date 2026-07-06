package editor

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

// extensionManagerImpl implements the ExtensionManager interface.
type extensionManagerImpl struct {
	log        zerolog.Logger
	extensions map[string]bool // simple in-memory state for scaffolding
	mu         sync.RWMutex
}

// NewExtensionManager creates a new instance of ExtensionManager.
func NewExtensionManager(logger zerolog.Logger) ExtensionManager {
	return &extensionManagerImpl{
		log:        logger.With().Str("component", "ExtensionManager").Logger(),
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

	// In a real implementation, this would dispatch a task to the Editor Provider
	// to execute `code-server --install-extension <extensionID>`
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

	// In a real implementation, this would dispatch a task to the Editor Provider
	// to execute `code-server --uninstall-extension <extensionID>`
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
