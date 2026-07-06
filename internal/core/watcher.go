package core

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

type WorkspaceWatcher struct {
	log      zerolog.Logger
	bus      events.Bus
	watcher  *fsnotify.Watcher
	roots    map[string]string
	mu       sync.Mutex
	debounce time.Duration
}

func NewWorkspaceWatcher(log zerolog.Logger, bus events.Bus) (*WorkspaceWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &WorkspaceWatcher{
		log:      log,
		bus:      bus,
		watcher:  w,
		roots:    make(map[string]string),
		debounce: 500 * time.Millisecond,
	}, nil
}

func (w *WorkspaceWatcher) AddWorkspace(id, root string) error {
	w.mu.Lock()
	w.roots[id] = root
	w.mu.Unlock()
	return w.watchRecursive(root)
}

// RemoveWorkspace stops watching the given workspace.
func (w *WorkspaceWatcher) RemoveWorkspace(id string) {
	w.mu.Lock()
	root, ok := w.roots[id]
	if ok {
		delete(w.roots, id)
	}
	w.mu.Unlock()
	if ok {
		_ = w.watcher.Remove(root)
	}
}

func (w *WorkspaceWatcher) watchRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return w.watcher.Add(path)
		}
		return nil
	})
}

func (w *WorkspaceWatcher) Run(ctx context.Context) error {
	defer w.watcher.Close()

	pending := make(map[string]time.Time)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	w.log.Info().Msg("workspace watcher started")

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			w.mu.Lock()
			var matchedID string
			for id, root := range w.roots {
				if strings.HasPrefix(event.Name, root) {
					matchedID = id
					break
				}
			}
			w.mu.Unlock()

			if matchedID != "" {
				if !strings.Contains(event.Name, ".devserver") && !strings.Contains(event.Name, ".git") {
					if event.Has(fsnotify.Create) {
						info, err := os.Stat(event.Name)
						if err == nil && info.IsDir() {
							_ = w.watchRecursive(event.Name)
						}
					}
					pending[matchedID] = time.Now()
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			w.log.Error().Err(err).Msg("watcher error")
		case <-ticker.C:
			now := time.Now()
			for id, t := range pending {
				if now.Sub(t) >= w.debounce {
					w.log.Debug().Str("workspace", id).Msg("publishing workspace change event")
					w.bus.Publish(events.WorkspaceChanged, events.WorkspaceChangedEvent{WorkspaceID: id})
					delete(pending, id)
				}
			}
		}
	}
}
