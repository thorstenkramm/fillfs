// Package app orchestrates the fillfs command execution.
package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/options"
	"github.com/thorstenkramm/fillfs/internal/plan"
	"github.com/thorstenkramm/fillfs/internal/registry"
	"github.com/thorstenkramm/fillfs/internal/runerr"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// Run executes fillfs with the provided config.
//
//nolint:funlen
func Run(ctx context.Context, cfg options.Config) error {
	gens := registry.Generators()
	fmt.Println("Generating plan...")
	p, err := plan.Build(cfg, gens)
	if err != nil {
		return fmt.Errorf("build plan: %w", err)
	}

	if err := ensureDisk(cfg.Dest, p.TotalSize); err != nil {
		return fmt.Errorf("check disk space: %w", err)
	}

	printSummary(cfg, p)
	if !cfg.Yes {
		ok, err := promptYes()
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Aborted.")
			return nil
		}
	}

	cachePath := cfg.CacheDir
	if !cfg.CacheIsDefault {
		cachePath = filepath.Join(cachePath, "fillfs")
	}

	cacheMgr := cache.New(cachePath, cfg.CacheIsDefault)
	if err := cacheMgr.Prepare(); err != nil {
		return fmt.Errorf("prepare cache: %w", err)
	}
	if cfg.CleanCache {
		defer func() {
			fmt.Println("Cleaning cache directory...")
			if err := cacheMgr.Clean(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to clean cache: %v\n", err)
			}
		}()
	}

	if err := prepareDestination(cfg); err != nil {
		return fmt.Errorf("prepare destination: %w", err)
	}

	fmt.Println("Creating directories...")
	for _, dir := range p.Directories {
		if err := os.MkdirAll(filepath.Join(cfg.Dest, dir.Path), 0o750); err != nil {
			return fmt.Errorf("create dir %s: %w", dir.Path, err)
		}
	}

	fmt.Println("Copying files...")
	genMap := mapGenerators(gens)
	for _, f := range p.Files {
		g, ok := genMap[f.Ext]
		if !ok {
			return fmt.Errorf("missing generator for %s", f.Ext)
		}

		seed := sources.Seed{URL: f.SeedURL, FileName: f.SeedName, Size: f.SeedSize, Extension: f.Ext}
		destPath := filepath.Join(cfg.Dest, f.DestPath)
		fmt.Printf("copy %s -> %s\n", seed.FileName, destPath)
		if err := g.Copy(ctx, cacheMgr, seed, destPath); err != nil {
			return fmt.Errorf("copy %s: %w", destPath, err)
		}
	}

	fmt.Println("Done.")
	return nil
}

func prepareDestination(cfg options.Config) error {
	info, err := os.Stat(cfg.Dest)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.Dest, 0o750); err != nil {
				return fmt.Errorf("create dest dir: %w", err)
			}
			return nil
		}
		return fmt.Errorf("stat dest: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("destination %s is not a directory", cfg.Dest)
	}

	empty, err := isDirEmpty(cfg.Dest)
	if err != nil {
		return fmt.Errorf("check dest emptiness: %w", err)
	}
	if empty {
		return nil
	}

	if !cfg.WipeDest {
		return runerr.WithCode(errors.New("destination is not empty"), 5) //nolint:wrapcheck
	}

	entries, err := os.ReadDir(cfg.Dest)
	if err != nil {
		return fmt.Errorf("read dest: %w", err)
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(cfg.Dest, entry.Name())); err != nil {
			return fmt.Errorf("remove existing file: %w", err)
		}
	}
	return nil
}

func isDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("read dir: %w", err)
	}
	return len(entries) == 0, nil
}

func mapGenerators(gens []generator.Generator) map[string]generator.Generator {
	m := make(map[string]generator.Generator, len(gens))
	for _, g := range gens {
		m[g.Extension()] = g
	}
	return m
}

func promptYes() (bool, error) {
	fmt.Print("Proceed? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
		return false, fmt.Errorf("read input: %w", err)
	}
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
}

func printSummary(cfg options.Config, p plan.Plan) {
	fmt.Println("Plan summary:")
	fmt.Printf("- Dest: %s\n", cfg.Dest)
	fmt.Printf("- Cache: %s\n", cfg.CacheDir)
	fmt.Printf("- Directories: %d\n", len(p.Directories))
	fmt.Printf("- Files: %d\n", len(p.Files))
	fmt.Printf("- Estimated size: %s\n", humanSize(p.TotalSize))
	fmt.Println("- Per extension:")
	for ext, count := range p.PerExtension {
		fmt.Printf("  %s: %d\n", ext, count)
	}
}

func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	n := float64(b)
	exp := 0
	for n >= unit {
		n /= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", n, "KMGTPE"[exp-1])
}

func ensureDisk(path string, required int64) error {
	if required <= 0 {
		return nil
	}

	target := path
	if target == "" {
		target = "."
	}

	if _, err := os.Stat(target); err != nil {
		if os.IsNotExist(err) {
			target = filepath.Dir(target)
		} else {
			return fmt.Errorf("stat target: %w", err)
		}
		if target == "" {
			target = "."
		}
	}

	var stat unix.Statfs_t
	if err := unix.Statfs(target, &stat); err != nil {
		return fmt.Errorf("statfs: %w", err)
	}

	blocks := stat.Bavail
	blockSize := uint64(stat.Bsize)
	var available uint64
	if blockSize > 0 && blocks > math.MaxUint64/blockSize {
		available = math.MaxUint64
	} else {
		available = blocks * blockSize
	}
	availableSigned := int64(math.Min(float64(available), float64(math.MaxInt64)))

	if required > availableSigned {
		return fmt.Errorf("disk space: %w", runerr.WithCode(fmt.Errorf(
			"not enough disk space: need %s, available %s",
			humanSize(required),
			humanSize(availableSigned),
		), 3))
	}

	return nil
}
