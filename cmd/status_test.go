package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/layout"
	"github.com/stretchr/testify/require"
)

func TestStatusCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))
		_, err := statusCmd(t)
		require.EqualError(t, err, layout.ErrNotTracRepository.Error())
	})

	t.Run("empty repository", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		out, err := statusCmd(t)
		require.NoError(t, err)
		require.Equal(t, "nothing to commit (create/copy files and use \"trac add\" to track)\n", out)
	})

	t.Run("untracked files", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		untrackedFilePath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(untrackedFilePath, []byte("some content"), 0644))

		out, err := statusCmd(t)
		require.NoError(t, err)
		require.Contains(t, out, "Untracked files:")
		require.Contains(t, out, "test.txt")
	})

	t.Run("tracked files", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		trackedFilePath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(trackedFilePath, []byte("tracked content"), 0644))
		require.NoError(t, addCmd(t, trackedFilePath))

		out, err := statusCmd(t)
		require.NoError(t, err)
		require.Contains(t, out, "Changes to be committed:")
		require.Contains(t, out, "new file:   test.txt")
	})
}
