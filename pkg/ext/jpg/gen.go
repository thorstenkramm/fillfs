// Package jpg generates .jpg files from seeds.
package jpg

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/copier"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// New returns a generator for .jpg files.
func New() generator.Generator { return gen{} }

type gen struct{}

func (gen) Extension() string { return ".jpg" }

func (gen) Seeds() []sources.Seed { return sources.SeedsByExtension(".jpg") }

func (gen) Copy(ctx context.Context, cacheMgr cache.Manager, seed sources.Seed, destPath string) error {
	return copier.Copy(ctx, cacheMgr, seed, destPath) //nolint:wrapcheck
}
