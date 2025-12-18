// Package options parses CLI flags and produces runtime configuration.
package options

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds runtime configuration parsed from flags.
type Config struct {
	Dest           string
	CacheDir       string
	CacheIsDefault bool
	CleanCache     bool
	Folders        int
	FilesPerFolder int
	Depths         float64
	Yes            bool
	WipeDest       bool
}

// Load parses CLI flags via viper/pflag and returns a validated Config.
func Load() (Config, error) {
	pflag.String("dest", ".", "Destination directory to fill")
	pflag.String("cache-dir", cacheDefault(), "Directory to cache seed files")
	pflag.Bool("clean-cache", false, "Remove cache directory before running")
	pflag.Int("folders", 2, "Number of folders to create per level")
	pflag.Int("files-per-folder", 20, "Number of files to create in each folder")
	pflag.Float64("depths", 1, "Depth of recursion (floats allowed)")
	pflag.Bool("yes", false, "Do not prompt for confirmation")
	pflag.Bool("wipe-dest", false, "Delete destination contents before filling")

	pflag.Parse()

	_ = viper.BindPFlags(pflag.CommandLine)

	dest := filepath.Clean(viper.GetString("dest"))
	cache := filepath.Clean(viper.GetString("cache-dir"))
	cacheDefaultUsed := !pflag.Lookup("cache-dir").Changed || cache == "" || cache == "."
	if cacheDefaultUsed {
		cache = cacheDefault()
	}

	cfg := Config{
		Dest:           dest,
		CacheDir:       cache,
		CacheIsDefault: cacheDefaultUsed,
		CleanCache:     viper.GetBool("clean-cache"),
		Folders:        viper.GetInt("folders"),
		FilesPerFolder: viper.GetInt("files-per-folder"),
		Depths:         viper.GetFloat64("depths"),
		Yes:            viper.GetBool("yes"),
		WipeDest:       viper.GetBool("wipe-dest"),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Folders <= 0 {
		return fmt.Errorf("folders must be positive")
	}
	if c.FilesPerFolder <= 0 {
		return fmt.Errorf("files-per-folder must be positive")
	}
	if c.Depths <= 0 {
		return fmt.Errorf("depths must be positive")
	}
	return nil
}

func cacheDefault() string {
	tmp := os.TempDir()
	if tmp == "" {
		tmp = "/tmp"
	}
	return filepath.Join(tmp, ".fillfs")
}
