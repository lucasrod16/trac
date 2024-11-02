package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, ".trac")

	cmd := NewInitCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{tmpDir})
	require.NoError(t, cmd.Execute())

	expectedDirs := []string{
		filepath.Join(repoPath, "objects"),
		filepath.Join(repoPath, "refs", "heads"),
	}
	for _, dir := range expectedDirs {
		_, err := os.Stat(dir)
		require.NoError(t, err, "Expected directory %s to exist, but it does not.", dir)
	}

	headPath := filepath.Join(repoPath, "HEAD")
	_, err := os.Stat(headPath)
	require.NoError(t, err, "Expected HEAD file to exist, but it does not.")

	mainBranchPath := filepath.Join(repoPath, "refs", "heads", "main")
	_, err = os.Stat(mainBranchPath)
	require.NoError(t, err, "Expected main branch reference file %s to exist, but it does not.", mainBranchPath)

	require.Equal(t, fmt.Sprintf("Initialized empty trac repository in %s\n", repoPath), buf.String(), "Expected initialization message not found")

	// verify re-init
	buf.Reset()
	require.NoError(t, cmd.Execute())
	require.Equal(t, fmt.Sprintf("Reinitialized existing trac repository in %s\n", repoPath), buf.String(), "Expected reiinitialization message not found")
}
