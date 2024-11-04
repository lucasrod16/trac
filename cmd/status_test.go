package cmd

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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
		_ = initRepository(t)

		cmd := NewStatusCmd()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		require.NoError(t, cmd.Execute())
		require.Equal(t, "nothing to commit (create/copy files and use \"trac add\" to track)\n", buf.String())
	})

	t.Run("untracked files", func(t *testing.T) {
		tmpdir := initRepository(t)

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

		repoPath := filepath.Join(tmpdir, ".trac")
		trackedFilePath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(trackedFilePath, []byte("tracked content"), 0644))

		// simulate staging a file by creating an object file in the .trac/objects directory
		content, err := os.ReadFile(trackedFilePath)
		require.NoError(t, err)
		hash := fmt.Sprintf("%x", sha256.Sum256(content))
		objectPath := filepath.Join(repoPath, "objects", hash[:2], hash[2:])
		require.NoError(t, os.MkdirAll(filepath.Dir(objectPath), 0755))
		require.NoError(t, os.WriteFile(objectPath, content, 0644))

		cmd := NewStatusCmd()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)

		require.NoError(t, cmd.Execute())
		require.Contains(t, buf.String(), "Changes to be committed:")
		require.Contains(t, buf.String(), "new file:   test.txt")
	})
}

// initRepository initializes a new trac repository for testing.
func initRepository(t *testing.T) string {
	t.Helper()

	tmpdir := t.TempDir()
	require.NoError(t, os.Chdir(tmpdir))

	initCmd := NewInitCmd()
	initCmd.SetOut(io.Discard)
	initCmd.SetErr(io.Discard)
	require.NoError(t, initCmd.Execute())

	return tmpdir
}
