// Package mp3 generates .mp3 files from seeds.
package mp3

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/copier"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// New returns a generator for .mp3 files.
func New() generator.Generator { return gen{} }

type gen struct{}

func (gen) Extension() string { return ".mp3" }

func (gen) Seeds() []sources.Seed { return sources.SeedsByExtension(".mp3") }

func (gen) Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	return copier.Copy(ctx, cacheMgr, seed, destPath) //nolint:wrapcheck
}
