package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasrod16/trac/internal/layout"
	"github.com/lucasrod16/trac/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAddCommand(t *testing.T) {
	t.Run("non-trac repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		require.NoError(t, os.Chdir(tmpdir))
		require.EqualError(t, addCmd(t, "."), layout.ErrNotTracRepository.Error())
	})

	t.Run("add file outside repository should error", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		outsidePath := filepath.Join(filepath.Dir(tmpdir), "test.txt")
		require.NoError(t, os.WriteFile(outsidePath, []byte("content"), 0644))

		err := addCmd(t, outsidePath)
		require.ErrorContains(t, err, "failed to add file")
		require.ErrorContains(t, err, "outside the repository")
	})

	t.Run("add a single file", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		testPath := filepath.Join(tmpdir, "test.txt")
		require.NoError(t, os.WriteFile(testPath, []byte("some content"), 0644))
		require.NoError(t, addCmd(t, testPath))

		idx := getIndex(t, tmpdir)

		expected, err := utils.HashFile(testPath)
		require.NoError(t, err)
		actual := idx.Staged[testPath]
		require.Equal(t, expected, actual)
	})

	t.Run("add multiple files", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		testPath1 := filepath.Join(tmpdir, "test1.txt")
		require.NoError(t, os.WriteFile(testPath1, []byte("content 1"), 0644))
		testPath2 := filepath.Join(tmpdir, "test2.txt")
		require.NoError(t, os.WriteFile(testPath2, []byte("content 2"), 0644))

		require.NoError(t, addCmd(t, testPath1, testPath2))

		idx := getIndex(t, tmpdir)

		expected, err := utils.HashFile(testPath1)
		require.NoError(t, err)
		actual := idx.Staged[testPath1]
		require.Equal(t, expected, actual)

		expected, err = utils.HashFile(testPath2)
		require.NoError(t, err)
		actual = idx.Staged[testPath2]
		require.Equal(t, expected, actual)
	})

	t.Run("add files recursively (working dir)", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		subdir := filepath.Join(tmpdir, "subdir")
		require.NoError(t, os.Mkdir(subdir, 0755))

		testPath1 := filepath.Join(subdir, "test1.txt")
		require.NoError(t, os.WriteFile(testPath1, []byte("content 1"), 0644))
		testPath2 := filepath.Join(subdir, "test2.txt")
		require.NoError(t, os.WriteFile(testPath2, []byte("content 2"), 0644))

		require.NoError(t, addCmd(t, "."))

		idx := getIndex(t, tmpdir)

		testPath1, err := filepath.Rel(tmpdir, testPath1)
		require.NoError(t, err)
		testPath2, err = filepath.Rel(tmpdir, testPath2)
		require.NoError(t, err)

		expected, err := utils.HashFile(testPath1)
		require.NoError(t, err)
		actual := idx.Staged[testPath1]
		require.Equal(t, expected, actual)

		expected, err = utils.HashFile(testPath2)
		require.NoError(t, err)
		actual = idx.Staged[testPath2]
		require.Equal(t, expected, actual)
	})

	t.Run("add files recursively", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))

		subdir := filepath.Join(tmpdir, "subdir")
		require.NoError(t, os.Mkdir(subdir, 0755))

		testPath1 := filepath.Join(subdir, "test1.txt")
		require.NoError(t, os.WriteFile(testPath1, []byte("content 1"), 0644))
		testPath2 := filepath.Join(subdir, "test2.txt")
		require.NoError(t, os.WriteFile(testPath2, []byte("content 2"), 0644))

		require.NoError(t, addCmd(t, subdir))

		idx := getIndex(t, tmpdir)

		testPath1, err := filepath.Rel(tmpdir, testPath1)
		require.NoError(t, err)
		testPath2, err = filepath.Rel(tmpdir, testPath2)
		require.NoError(t, err)

		expected, err := utils.HashFile(testPath1)
		require.NoError(t, err)
		actual := idx.Staged[testPath1]
		require.Equal(t, expected, actual)

		expected, err = utils.HashFile(testPath2)
		require.NoError(t, err)
		actual = idx.Staged[testPath2]
		require.Equal(t, expected, actual)
	})

	t.Run("invalid file path", func(t *testing.T) {
		tmpdir := initRepository(t)
		require.NoError(t, os.Chdir(tmpdir))
		invalidFilePath := filepath.Join(tmpdir, "invalid.txt")
		require.Error(t, addCmd(t, invalidFilePath))
	})
}
