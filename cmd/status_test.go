package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatusCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))

		cmd := NewStatusCmd()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		require.EqualError(t, cmd.Execute(), "not a trac repository (or any of the parent directories): .trac")
	})

	t.Run("empty repository", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		cmd := NewStatusCmd()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		require.NoError(t, cmd.Execute())
		require.Equal(t, "nothing to commit (create/copy files and use \"trac add\" to track)\n", buf.String())
	})

	t.Run("untracked files", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		untrackedFilePath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(untrackedFilePath, []byte("some content"), 0644))

		cmd := NewStatusCmd()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		require.NoError(t, cmd.Execute())
		require.Contains(t, buf.String(), "Untracked files:")
		require.Contains(t, buf.String(), "test.txt")
	})

	t.Run("tracked files", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		trackedFilePath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(trackedFilePath, []byte("tracked content"), 0644))

		addCmd := NewAddCmd()
		addCmd.SetArgs([]string{trackedFilePath})
		addCmd.SetOut(io.Discard)
		addCmd.SetErr(io.Discard)
		require.NoError(t, addCmd.Execute())

		cmd := NewStatusCmd()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		require.NoError(t, cmd.Execute())
		require.Contains(t, buf.String(), "Changes to be committed:")
		require.Contains(t, buf.String(), "new file:   test.txt")
	})
}
