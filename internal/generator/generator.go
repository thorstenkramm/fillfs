// Package generator defines the extension-specific generator interface.
package generator

import (
	"context"

	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// Generator creates files for a specific extension.
type Generator interface {
	Extension() string
	Seeds() []sources.Seed
	Copy(ctx context.Context, cache cache.Manager, seed sources.Seed, destPath string) error
}
