// Package filesystem provides a pluggable abstraction over filesystem
// operations. Today only a local implementation exists. Future
// implementations will support SSH, Docker, S3, and remote agents
// without changing any consumer code.
package filesystem

import (
	"io"
	"time"
)

// FileInfo describes a single entry in the filesystem.
type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
	IsDir   bool
	Mode    uint32
}

// FS is the core interface every WorkspaceManager, Indexer, and future
// subsystem should accept instead of using os.* directly.
//
// Implementations must be safe for concurrent use.
type FS interface {
	// Read returns the full contents of the file at path.
	Read(path string) ([]byte, error)

	// Write atomically writes data to path, creating parent directories
	// as needed.
	Write(path string, data []byte) error

	// Delete removes the file or empty directory at path.
	Delete(path string) error

	// Exists reports whether path exists.
	Exists(path string) (bool, error)

	// Stat returns metadata for the path.
	Stat(path string) (FileInfo, error)

	// ReadDir lists the immediate children of dir.
	ReadDir(dir string) ([]FileInfo, error)

	// Walk calls fn for every file and directory under root, respecting
	// skip semantics (return filepath.SkipDir to skip a subtree).
	Walk(root string, fn WalkFunc) error

	// Open returns a ReadCloser for streaming large files.
	Open(path string) (io.ReadCloser, error)

	// MkdirAll creates dir and all parent directories.
	MkdirAll(dir string) error
}

// WalkFunc is the callback signature for FS.Walk.
type WalkFunc func(path string, info FileInfo, err error) error
