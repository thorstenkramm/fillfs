// Package ppt generates .ppt files from seeds.
package ppt

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/copier"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// New returns a generator for .ppt files.
func New() generator.Generator { return gen{} }

type gen struct{}

func (gen) Extension() string { return ".ppt" }

func (gen) Seeds() []sources.Seed { return sources.SeedsByExtension(".ppt") }

func (gen) Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	return copier.Copy(ctx, cacheMgr, seed, destPath) //nolint:wrapcheck
}
