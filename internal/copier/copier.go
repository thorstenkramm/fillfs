// Package copier handles copying cached seeds into destination files.
package copier

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// Copy downloads the seed into cache if necessary and copies it to destPath.
func Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	srcPath, err := cacheMgr.Ensure(ctx, seed)
	if err != nil {
		return fmt.Errorf("ensure cache for %s: %w", seed.FileName, err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
		return fmt.Errorf("mkdir for %s: %w", destPath, err)
	}

	src, err := os.Open(srcPath) //nolint:gosec // path comes from controlled cache
	if err != nil {
		return fmt.Errorf("open source %s: %w", srcPath, err)
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(destPath) //nolint:gosec // destination is intended by tool
	if err != nil {
		return fmt.Errorf("create dest %s: %w", destPath, err)
	}
	defer func() {
		_ = dst.Close()
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy to %s: %w", destPath, err)
	}

	return nil
}
