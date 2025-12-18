// Package xlsx generates .xlsx files from seeds.
package xlsx

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/copier"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// New returns a generator for .xlsx files.
func New() generator.Generator { return gen{} }

type gen struct{}

func (gen) Extension() string { return ".xlsx" }

func (gen) Seeds() []sources.Seed { return sources.SeedsByExtension(".xlsx") }

func (gen) Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	return copier.Copy(ctx, cacheMgr, seed, destPath) //nolint:wrapcheck
}
