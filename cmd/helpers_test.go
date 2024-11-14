package cmd

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
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

func addCmd(t *testing.T, args ...string) error {
	t.Helper()
	cmd := NewAddCmd()
	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd.Execute()
}

func commitCmd(t *testing.T, args ...string) error {
	t.Helper()
	cmd := NewCommitCmd()
	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd.Execute()
}

func statusCmd(t *testing.T) (output string, err error) {
	t.Helper()
	cmd := NewStatusCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err = cmd.Execute()
	return buf.String(), err
}

func getIndex(t *testing.T, repoPath string) *index.Index {
	t.Helper()
	l, err := layout.New(repoPath)
	require.NoError(t, err)
	idx := index.New()
	require.NoError(t, idx.Load(l))
	return idx
}
