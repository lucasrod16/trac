package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCommitCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))

		cmd := NewCommitCmd()
		cmd.SetArgs([]string{"-m", "test commit message"})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		require.EqualError(t, cmd.Execute(), "not a trac repository (or any of the parent directories): .trac")
	})

	t.Run("verify commit writes to object database", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		testFile := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))

		addCmd := NewAddCmd()
		addCmd.SetArgs([]string{testFile})
		addCmd.SetOut(io.Discard)
		addCmd.SetErr(io.Discard)
		require.NoError(t, addCmd.Execute())

		commitCmd := NewCommitCmd()
		commitCmd.SetArgs([]string{"-m", "test commit message"})
		commitCmd.SetOut(io.Discard)
		commitCmd.SetErr(io.Discard)
		require.NoError(t, commitCmd.Execute())

		hash, err := utils.HashFile(testFile)
		require.NoError(t, err)
		require.FileExists(t, filepath.Join(".trac", "objects", hash[:2], hash[2:]))
	})
}
