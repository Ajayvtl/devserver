package core

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/filesystem"
	"github.com/rs/zerolog"
)

// WorkspaceEntry holds the runtime state of a single managed workspace.
type WorkspaceEntry struct {
	ID     string
	Root   string
	Active bool
}

// WorkspaceManager is the single owner of workspace lifecycle.
// It coordinates the Indexer, Watcher, and Provider so that no
// other code needs to call those subsystems directly for workspace
// registration or teardown.
type WorkspaceManager struct {
	log     zerolog.Logger
	mu      sync.RWMutex
	entries map[string]*WorkspaceEntry

	indexer  *Indexer
	watcher  *WorkspaceWatcher
	provider *WorkspaceProvider
	bus      events.Bus
	fs       filesystem.FS
}

// WorkspaceManagerDeps collects all dependencies for the manager.
type WorkspaceManagerDeps struct {
	Logger   zerolog.Logger
	Indexer  *Indexer
	Watcher  *WorkspaceWatcher
	Provider *WorkspaceProvider
	Bus      events.Bus
	FS       filesystem.FS
}

// NewWorkspaceManager creates a manager that coordinates all workspace
// subsystems through a single API.
func NewWorkspaceManager(deps WorkspaceManagerDeps) *WorkspaceManager {
	return &WorkspaceManager{
		log:      deps.Logger,
		entries:  make(map[string]*WorkspaceEntry),
		indexer:  deps.Indexer,
		watcher:  deps.Watcher,
		provider: deps.Provider,
		bus:      deps.Bus,
		fs:       deps.FS,
	}
}

// Open registers a workspace across every subsystem and runs the
// first synchronous index so the workspace is immediately queryable.
func (m *WorkspaceManager) Open(ctx context.Context, id, root string) error {
	m.mu.Lock()
	if _, exists := m.entries[id]; exists {
		m.mu.Unlock()
		return fmt.Errorf("workspace %q is already open", id)
	}
	m.entries[id] = &WorkspaceEntry{ID: id, Root: root, Active: true}
	m.mu.Unlock()

	// Verify the root exists on the filesystem.
	exists, err := m.fs.Exists(root)
	if err != nil {
		m.removeEntry(id)
		return fmt.Errorf("check workspace root: %w", err)
	}
	if !exists {
		m.removeEntry(id)
		return fmt.Errorf("workspace root %q does not exist", root)
	}

	// Register with Indexer.
	m.indexer.AddWorkspace(id, root)

	// Register with Watcher.
	if m.watcher != nil {
		if err := m.watcher.AddWorkspace(id, root); err != nil {
			m.log.Warn().Err(err).Str("workspace", id).Msg("watcher registration failed, continuing without live watch")
		}
	}

	// Run the first synchronous index so data is available immediately.
	if err := m.indexer.Refresh(ctx, id); err != nil {
		m.log.Warn().Err(err).Str("workspace", id).Msg("initial index failed")
	}

	// Invalidate provider cache to pick up new data.
	if m.provider != nil {
		m.provider.Invalidate(id)
	}

	m.bus.Publish(events.WorkspaceActivated, events.WorkspaceActivatedEvent{
		WorkspaceID: id,
		Root:        root,
	})

	m.log.Info().Str("workspace", id).Str("root", root).Msg("workspace opened")
	return nil
}

// Close tears down a workspace from every subsystem.
func (m *WorkspaceManager) Close(id string) error {
	m.mu.Lock()
	entry, exists := m.entries[id]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("workspace %q is not open", id)
	}
	entry.Active = false
	delete(m.entries, id)
	m.mu.Unlock()

	// Tear down in reverse order of registration.
	if m.watcher != nil {
		m.watcher.RemoveWorkspace(id)
	}
	m.indexer.RemoveWorkspace(id)
	if m.provider != nil {
		m.provider.Invalidate(id)
	}

	m.bus.Publish(events.WorkspaceDeactivated, events.WorkspaceDeactivatedEvent{
		WorkspaceID: id,
	})

	m.log.Info().Str("workspace", id).Msg("workspace closed")
	return nil
}

// Get returns the workspace entry if it exists.
func (m *WorkspaceManager) Get(id string) (*WorkspaceEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[id]
	return e, ok
}

// List returns all active workspace IDs, sorted.
func (m *WorkspaceManager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]string, 0, len(m.entries))
	for id := range m.entries {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Validate checks whether the workspace root still exists on disk.
func (m *WorkspaceManager) Validate(id string) error {
	m.mu.RLock()
	entry, ok := m.entries[id]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("workspace %q not found", id)
	}
	exists, err := m.fs.Exists(entry.Root)
	if err != nil {
		return fmt.Errorf("validate workspace root: %w", err)
	}
	if !exists {
		return fmt.Errorf("workspace root %q no longer exists", entry.Root)
	}
	return nil
}

func (m *WorkspaceManager) removeEntry(id string) {
	m.mu.Lock()
	delete(m.entries, id)
	m.mu.Unlock()
}
