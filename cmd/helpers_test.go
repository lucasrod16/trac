package cmd

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// initRepository initializes a new trac repository for testing.
func initRepository(t *testing.T) string {
	t.Helper()

	tmpdir := t.TempDir()
	tmpdir, err := filepath.EvalSymlinks(tmpdir)
	require.NoError(t, err)

	initCmd := NewInitCmd()
	initCmd.SetOut(io.Discard)
	initCmd.SetErr(io.Discard)
	initCmd.SetArgs([]string{tmpdir})
	require.NoError(t, initCmd.Execute())

	return tmpdir
}
