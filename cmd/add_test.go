package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/index"
	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAddCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))

		cmd := NewAddCmd()
		cmd.SetArgs([]string{"."})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		require.EqualError(t, cmd.Execute(), "not a trac repository (or any of the parent directories): .trac")
	})

	t.Run("add a single file", func(t *testing.T) {
		tmpdir := initRepository(t)

		testPath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(testPath, []byte("some content"), 0644))

		cmd := NewAddCmd()
		cmd.SetArgs([]string{testPath})
		require.NoError(t, cmd.Execute())

		l, err := layout.New(tmpdir)
		require.NoError(t, err)
		idx := index.New()
		require.NoError(t, idx.Load(l))

		hash, err := utils.HashFile(testPath)
		require.NoError(t, err)
		require.Contains(t, idx.Staged, hash)
	})

	t.Run("add multiple files", func(t *testing.T) {
		tmpdir := initRepository(t)

		testPath1 := filepath.Join(tmpdir, "test1.txt")
		require.NoError(t, os.WriteFile(testPath1, []byte("content 1"), 0644))
		testPath2 := filepath.Join(tmpdir, "test2.txt")
		require.NoError(t, os.WriteFile(testPath2, []byte("content 2"), 0644))

		cmd := NewAddCmd()
		cmd.SetArgs([]string{testPath1, testPath2})
		require.NoError(t, cmd.Execute())

		l, err := layout.New(tmpdir)
		require.NoError(t, err)
		idx := index.New()
		require.NoError(t, idx.Load(l))

		hash1, err := utils.HashFile(testPath1)
		require.NoError(t, err)
		hash2, err := utils.HashFile(testPath2)
		require.NoError(t, err)

		require.Contains(t, idx.Staged, hash1)
		require.Contains(t, idx.Staged, hash2)
	})

	t.Run("add files recursively", func(t *testing.T) {
		tmpdir := initRepository(t)

		subdir := filepath.Join(tmpdir, "subdir")
		require.NoError(t, os.Mkdir(subdir, 0755))

		testPath1 := filepath.Join(subdir, "test1.txt")
		require.NoError(t, os.WriteFile(testPath1, []byte("content 1"), 0644))
		testPath2 := filepath.Join(subdir, "test2.txt")
		require.NoError(t, os.WriteFile(testPath2, []byte("content 2"), 0644))

		cmd := NewAddCmd()
		cmd.SetArgs([]string{"."})
		require.NoError(t, cmd.Execute())

		l, err := layout.New(tmpdir)
		require.NoError(t, err)
		idx := index.New()
		require.NoError(t, idx.Load(l))

		hash1, err := utils.HashFile(testPath1)
		require.NoError(t, err)
		hash2, err := utils.HashFile(testPath2)
		require.NoError(t, err)

		require.Contains(t, idx.Staged, hash1)
		require.Contains(t, idx.Staged, hash2)
	})

	t.Run("invalid file path", func(t *testing.T) {
		tmpdir := initRepository(t)

		invalidFilePath := filepath.Join(tmpdir, "invalid.txt")
		cmd := NewAddCmd()
		cmd.SetArgs([]string{invalidFilePath})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		err := cmd.Execute()
		require.ErrorContains(t, err, "failed to add file")
	})
}
