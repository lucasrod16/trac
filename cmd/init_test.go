package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	tmpdir := t.TempDir()
	repoPath := filepath.Join(tmpdir, ".trac")

	cmd := NewInitCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{tmpdir})
	require.NoError(t, cmd.Execute())

	expectedDirs := []string{
		filepath.Join(repoPath, "objects"),
		filepath.Join(repoPath, "refs", "heads"),
	}
	for _, dir := range expectedDirs {
		require.DirExists(t, dir)
	}

	require.FileExists(t, filepath.Join(repoPath, "HEAD"))
	require.FileExists(t, filepath.Join(repoPath, "refs", "heads", "main"))
	require.Equal(t, fmt.Sprintf("Initialized empty trac repository in %s\n", repoPath), buf.String())

	// verify re-init
	buf.Reset()
	require.NoError(t, cmd.Execute())
	require.Equal(t, fmt.Sprintf("Reinitialized existing trac repository in %s\n", repoPath), buf.String())
}
