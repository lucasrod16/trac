package layout

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	tmpdir := t.TempDir()

	actual, err := New(tmpdir)
	require.NoError(t, err)
	require.NotNil(t, actual)

	expected := &Layout{
		Root:     tmpdir,
		Config:   filepath.Join(tmpdir, ".trac"),
		Objects:  filepath.Join(tmpdir, ".trac", "objects"),
		HeadFile: filepath.Join(tmpdir, ".trac", "HEAD"),
		Index:    filepath.Join(tmpdir, ".trac", "index.json"),
	}
	require.Equal(t, expected, actual)
}

func TestExists(t *testing.T) {
	t.Parallel()
	l, err := New(t.TempDir())
	require.NoError(t, err)
	require.False(t, l.Exists())
	require.NoError(t, l.Init())
	require.True(t, l.Exists())
}

func TestValidateIsRepo(t *testing.T) {
	t.Parallel()
	l, err := New(t.TempDir())
	require.NoError(t, err)
	require.Error(t, l.ValidateIsRepo())
	require.NoError(t, l.Init())
	require.NoError(t, l.ValidateIsRepo())
}
