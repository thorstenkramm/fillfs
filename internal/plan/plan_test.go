package plan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thorstenkramm/fillfs/internal/cache"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/options"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

type stubGen struct {
	ext   string
	seeds []sources.Seed
}

func (g stubGen) Extension() string { return g.ext }

func (g stubGen) Seeds() []sources.Seed { return g.seeds }

func (g stubGen) Copy(_ context.Context, _ cache.Manager, _ sources.Seed, _ string) error { return nil }

func TestBuildPlanCounts(t *testing.T) {
	cfg := options.Config{Folders: 2, FilesPerFolder: 3, Depths: 2, Dest: "/tmp/d", CacheDir: "/tmp/c"}

	genA := stubGen{ext: ".a", seeds: []sources.Seed{{FileName: "a", Size: 10, Extension: ".a", URL: "http://example/a"}}}
	genB := stubGen{ext: ".b", seeds: []sources.Seed{{FileName: "b", Size: 10, Extension: ".b", URL: "http://example/b"}}}

	p, err := Build(cfg, []generator.Generator{genA, genB})
	assert.NoError(t, err)
	assert.Len(t, p.Directories, 6)
	assert.Len(t, p.Files, 18)
	assert.Equal(t, int64(180), p.TotalSize)

	low, high := extremes(p.PerExtension)
	assert.LessOrEqual(t, high-low, 1)
}

func TestBuildPlanFractionalDepth(t *testing.T) {
	cfg := options.Config{Folders: 3, FilesPerFolder: 1, Depths: 1.5, Dest: "/tmp/d", CacheDir: "/tmp/c"}
	gen := stubGen{ext: ".x", seeds: []sources.Seed{{FileName: "x", Size: 1, Extension: ".x", URL: "http://example/x"}}}

	p, err := Build(cfg, []generator.Generator{gen})
	assert.NoError(t, err)

	// Depth 1.5 => first level 3 dirs, partial next level round(3*0.5)=2 per parent => 6 dirs
	assert.Len(t, p.Directories, 9)
	assert.Len(t, p.Files, 9)
}

func extremes(m map[string]int) (int, int) {
	minCount, maxCount := int(^uint(0)>>1), 0
	for _, v := range m {
		if v < minCount {
			minCount = v
		}
		if v > maxCount {
			maxCount = v
		}
	}
	return minCount, maxCount
}
