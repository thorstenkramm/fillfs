// Package plan builds the deterministic plan of directories and files to create.
package plan

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/thorstenkramm/fillfs/internal/filenames"
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/internal/options"
	"github.com/thorstenkramm/fillfs/internal/sources"
)

// DirectoryPlan represents a directory to create relative to destination.
type DirectoryPlan struct {
	Path string
}

// FilePlan represents a file copy to execute.
type FilePlan struct {
	DestPath string
	SeedName string
	SeedSize int64
	SeedURL  string
	Ext      string
}

// Plan holds the directories, files, and disk usage estimate.
type Plan struct {
	Directories  []DirectoryPlan
	Files        []FilePlan
	TotalSize    int64
	PerExtension map[string]int
}

// Build constructs a deterministic plan from the provided config and generators.
func Build(cfg options.Config, gens []generator.Generator) (Plan, error) {
	if len(gens) == 0 {
		return Plan{}, errors.New("no generators registered")
	}

	chooser := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // not security sensitive

	dirs := generateDirectories(cfg)

	files, perExt, totalSize, err := generateFiles(cfg, dirs, gens, chooser)
	if err != nil {
		return Plan{}, fmt.Errorf("generate files: %w", err)
	}

	return Plan{Directories: dirs, Files: files, TotalSize: totalSize, PerExtension: perExt}, nil
}

func generateDirectories(cfg options.Config) []DirectoryPlan {
	dInt := int(math.Floor(cfg.Depths))
	dFrac := cfg.Depths - float64(dInt)

	capacity := cfg.Folders * int(math.Ceil(cfg.Depths))
	if capacity < cfg.Folders {
		capacity = cfg.Folders
	}
	dirs := make([]DirectoryPlan, 0, capacity)

	topCount := cfg.Folders
	if dInt == 0 && dFrac > 0 {
		topCount = int(math.Round(float64(cfg.Folders) * dFrac))
		if topCount == 0 {
			topCount = 1
		}
	}

	currentLevel := make([]string, 0, topCount)
	currentLevel = append(currentLevel, generateUniqueDirectoryNames(topCount)...)
	for _, name := range currentLevel {
		dirs = append(dirs, DirectoryPlan{Path: name})
	}

	if dInt > 1 {
		for range make([]struct{}, dInt-1) {
			next := make([]string, 0, len(currentLevel)*cfg.Folders)
			for _, parent := range currentLevel {
				children := generateUniqueDirectoryNames(cfg.Folders)
				for _, child := range children {
					path := filepath.Join(parent, child)
					dirs = append(dirs, DirectoryPlan{Path: path})
					next = append(next, path)
				}
			}
			currentLevel = next
		}
	}

	if dFrac > 0 {
		partialCount := int(math.Round(float64(cfg.Folders) * dFrac))
		if partialCount > 0 {
			for _, parent := range currentLevel {
				children := generateUniqueDirectoryNames(partialCount)
				for _, child := range children {
					path := filepath.Join(parent, child)
					dirs = append(dirs, DirectoryPlan{Path: path})
				}
			}
		}
	}

	return dirs
}

func generateUniqueDirectoryNames(count int) []string {
	seen := make(map[string]struct{}, count)
	namesOut := make([]string, 0, count)
	for len(namesOut) < count {
		name := filenames.RandomDirectoryName()
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		namesOut = append(namesOut, name)
	}
	return namesOut
}

func generateFiles(
	cfg options.Config,
	dirs []DirectoryPlan,
	gens []generator.Generator,
	chooser *rand.Rand,
) ([]FilePlan, map[string]int, int64, error) {
	if len(dirs) == 0 {
		return nil, nil, 0, fmt.Errorf("no directories generated")
	}

	extGenerators := make(map[string]generator.Generator)
	extOrder := make([]string, 0, len(gens))
	for _, g := range gens {
		extGenerators[g.Extension()] = g
		extOrder = append(extOrder, g.Extension())
	}

	counts := make(map[string]int)
	files := make([]FilePlan, 0, len(dirs)*cfg.FilesPerFolder)
	var totalSize int64

	for _, dir := range dirs {
		usedNames := map[string]struct{}{}
		for i := 0; i < cfg.FilesPerFolder; i++ {
			ext := pickExtension(counts, extOrder, chooser)
			gen, ok := extGenerators[ext]
			if !ok {
				return nil, nil, 0, fmt.Errorf("no generator for extension %s", ext)
			}

			seed := pickSeed(gen, chooser)
			if seed.FileName == "" {
				return nil, nil, 0, fmt.Errorf("no seeds for extension %s", ext)
			}

			name := randomFileName(usedNames, ext)
			dest := filepath.Join(dir.Path, name)
			file := FilePlan{
				DestPath: dest,
				SeedName: seed.FileName,
				SeedSize: seed.Size,
				SeedURL:  seed.URL,
				Ext:      ext,
			}
			files = append(files, file)
			counts[ext]++
			totalSize += seed.Size
		}
	}

	return files, counts, totalSize, nil
}

func pickExtension(counts map[string]int, exts []string, rnd *rand.Rand) string {
	minCount := math.MaxInt
	var candidates []string
	for _, ext := range exts {
		c := counts[ext]
		if c < minCount {
			minCount = c
			candidates = []string{ext}
			continue
		}
		if c == minCount {
			candidates = append(candidates, ext)
		}
	}
	if len(candidates) == 0 {
		return exts[0]
	}
	return candidates[rnd.Intn(len(candidates))]
}

func pickSeed(gen generator.Generator, rnd *rand.Rand) sources.Seed {
	seeds := gen.Seeds()
	if len(seeds) == 0 {
		return sources.Seed{}
	}
	return seeds[rnd.Intn(len(seeds))]
}

func randomFileName(used map[string]struct{}, ext string) string {
	for {
		name := randomBaseNameForExt(ext) + ext
		if _, exists := used[name]; exists {
			continue
		}
		used[name] = struct{}{}
		return name
	}
}

func randomBaseNameForExt(ext string) string {
	switch ext {
	case ".doc", ".docx", ".pdf", ".rtf", ".odt":
		return filenames.RandomDocumentFileName()
	case ".ppt":
		return filenames.RandomPowerpointFileName()
	case ".xlsx":
		return filenames.RandomSpreadsheetFileName()
	case ".jpg", ".webp":
		return filenames.RandomImageFileName()
	case ".mp3", ".ogg":
		return filenames.RandomSoundFileName()
	default:
		return filenames.RandomDocumentFileName()
	}
}
