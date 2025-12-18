// Package cache manages download and reuse of seed files.
package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thorstenkramm/fillfs/internal/runerr"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

const markerName = ".fillfs"

// Manager handles cached seed downloads.
type Manager struct {
	path   string
	client *http.Client
	mark   bool
}

// New creates a cache manager rooted at path.
// If mark is true, a marker file is required/created to identify fillfs ownership.
func New(path string, mark bool) Manager {
	return Manager{path: path, client: &http.Client{}, mark: mark}
}

// Path returns the cache root.
func (m Manager) Path() string {
	return m.path
}

// Clean deletes the cache directory.
func (m Manager) Clean() error {
	if err := os.RemoveAll(m.path); err != nil {
		return fmt.Errorf("remove cache: %w", err)
	}
	return nil
}

// Prepare ensures the cache directory exists and is marked for fillfs use.
func (m Manager) Prepare() error {
	info, err := os.Stat(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(m.path, 0o750); err != nil {
				return fmt.Errorf("create cache dir: %w", err)
			}
			if m.mark {
				return m.markDir()
			}
			return nil
		}
		return fmt.Errorf("stat cache dir: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("cache path %s is not a directory", m.path)
	}

	if m.mark {
		marker := filepath.Join(m.path, markerName)
		if _, err := os.Stat(marker); err != nil {
			if os.IsNotExist(err) {
				entries, err := os.ReadDir(m.path)
				if err != nil {
					return fmt.Errorf("read cache dir: %w", err)
				}
				if len(entries) > 0 {
					return runerr.WithCode(errors.New("cache directory exists but is not marked as fillfs"), 4) //nolint:wrapcheck
				}
				return m.markDir()
			}
			return fmt.Errorf("stat cache marker: %w", err)
		}
	}

	return nil
}

func (m Manager) markDir() error {
	marker := filepath.Join(m.path, markerName)
	f, err := os.Create(marker) //nolint:gosec // marker path is controlled by tool
	if err != nil {
		return fmt.Errorf("create cache marker: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close cache marker: %w", err)
	}
	return nil
}

// Ensure returns the local path for a seed, downloading it if missing or wrong size.
func (m Manager) Ensure(ctx context.Context, seed sources.Seed) (string, error) {
	if err := m.Prepare(); err != nil {
		return "", fmt.Errorf("prepare cache: %w", err)
	}

	dest := filepath.Join(m.path, seed.FileName)
	if info, err := os.Stat(dest); err == nil {
		if info.Size() == seed.Size {
			return dest, nil
		}
		if err := os.Remove(dest); err != nil {
			return "", fmt.Errorf("remove stale cache file: %w", err)
		}
	}

	if err := m.download(ctx, seed.URL, dest); err != nil {
		return "", fmt.Errorf("download seed: %w", err)
	}

	return dest, nil
}

func (m Manager) download(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	tmp := dest + ".part"
	if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
		return fmt.Errorf("create cache parent: %w", err)
	}

	out, err := os.Create(tmp) //nolint:gosec // destination is controlled by config; creation is intended
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = out.Close()
		return fmt.Errorf("copy body: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmp, dest); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
