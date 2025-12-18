// Package rtf generates .rtf files from seeds.
package rtf

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/copier"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// New returns a generator for .rtf files.
func New() generator.Generator { return gen{} }

type gen struct{}

func (gen) Extension() string { return ".rtf" }

func (gen) Seeds() []sources.Seed { return sources.SeedsByExtension(".rtf") }

func (gen) Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	return copier.Copy(ctx, cacheMgr, seed, destPath) //nolint:wrapcheck
}
