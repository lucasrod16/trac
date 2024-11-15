package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/commit"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCommitCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))
		err := commitCmd(t, "-m", "test commit message")
		require.EqualError(t, err, layout.ErrNotTracRepository.Error())
	})

	t.Run("commit with no index file", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))
		err := commitCmd(t, "-m", "test commit message")
		require.EqualError(t, err, commit.ErrNothingAddedToCommit.Error())
	})

	t.Run("commit with clean working tree", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		testFile := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))
		require.NoError(t, addCmd(t, testFile))

		require.NoError(t, commitCmd(t, "-m", "test commit message"))
		err := commitCmd(t, "-m", "test commit message")
		require.EqualError(t, err, commit.ErrWorkingTreeClean.Error())
	})

	t.Run("verify commit writes to object database", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		testFile := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))
		require.NoError(t, addCmd(t, testFile))

		require.NoError(t, commitCmd(t, "-m", "test commit message"))

		hash, err := utils.HashFile(testFile)
		require.NoError(t, err)
		require.FileExists(t, filepath.Join(".trac", "objects", hash[:2], hash[2:]))
		latestCommit, err := os.ReadFile(filepath.Join(".trac", "HEAD"))
		require.NoError(t, err)
		require.NotEmpty(t, latestCommit)
	})
}
