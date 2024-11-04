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
		RootPath:      filepath.Join(tmpdir, ".trac"),
		ObjectsPath:   filepath.Join(tmpdir, ".trac", "objects"),
		RefsPath:      filepath.Join(tmpdir, ".trac", "refs"),
		HeadsPath:     filepath.Join(tmpdir, ".trac", "refs", "heads"),
		HeadFilePath:  filepath.Join(tmpdir, ".trac", "HEAD"),
		MainBranchRef: filepath.Join(tmpdir, ".trac", "refs", "heads", "main"),
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
