package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Local implements FS using the host operating system.
type Local struct{}

// NewLocal returns a local filesystem implementation.
func NewLocal() *Local {
	return &Local{}
}

func (l *Local) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (l *Local) Write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (l *Local) Delete(path string) error {
	return os.Remove(path)
}

func (l *Local) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (l *Local) Stat(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}
	return toFileInfo(path, info), nil
}

func (l *Local) ReadDir(dir string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]FileInfo, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, toFileInfo(filepath.Join(dir, e.Name()), info))
	}
	return out, nil
}

func (l *Local) Walk(root string, fn WalkFunc) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fn(path, FileInfo{}, err)
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return fn(path, FileInfo{Path: path, IsDir: d.IsDir()}, infoErr)
		}
		walkErr := fn(path, toFileInfo(path, info), nil)
		// Translate filepath.SkipDir so callers can use it.
		if walkErr == filepath.SkipDir {
			return filepath.SkipDir
		}
		return walkErr
	})
}

func (l *Local) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *Local) MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func toFileInfo(path string, info os.FileInfo) FileInfo {
	return FileInfo{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime().UTC(),
		IsDir:   info.IsDir(),
		Mode:    uint32(info.Mode()),
	}
}
