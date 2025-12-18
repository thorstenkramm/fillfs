package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thorstenkramm/fillfs/internal/runerr"
)

func TestPrepareCreatesMarker(t *testing.T) {
	dir := t.TempDir()
	mgr := New(filepath.Join(dir, "cache"), true)

	assert.NoError(t, mgr.Prepare())
	_, err := os.Stat(filepath.Join(mgr.Path(), markerName))
	assert.NoError(t, err)
}

func TestPrepareRejectsForeignContent(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "cache")
	assert.NoError(t, os.MkdirAll(cacheDir, 0o750))
	assert.NoError(t, os.WriteFile(filepath.Join(cacheDir, "foreign"), []byte("x"), 0o600))

	mgr := New(cacheDir, true)
	err := mgr.Prepare()
	assert.Error(t, err)
	assert.Equal(t, 4, runerr.Code(err, 1))
}
