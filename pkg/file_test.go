package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/kilianpaquier/filesystem/pkg/tests"
)

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "file.txt")
	dest := filepath.Join(tmp, "copy.txt")

	err := os.WriteFile(src, []byte("hey file"), filesystem.RwRR)
	require.NoError(t, err)

	t.Run("error_src_not_exists", func(t *testing.T) {
		// Arrange
		src := filepath.Join(tmp, "invalid.txt")

		// Act
		err := filesystem.CopyFile(src, dest)

		// Assert
		assert.ErrorContains(t, err, "failed to read")
		assert.NoFileExists(t, dest)
	})

	t.Run("error_destdir_not_exists", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(tmp, "invalid", "file.txt")

		// Act
		err := filesystem.CopyFile(src, dest)

		// Assert
		assert.ErrorContains(t, err, "failed to create")
		assert.NoFileExists(t, dest)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		err := filesystem.CopyFile(src, dest)

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})

	t.Run("success_with_fs", func(t *testing.T) {
		// Act
		err := filesystem.CopyFile(src, dest,
			filesystem.WithFS(filesystem.OS()),
			filesystem.WithJoin(filepath.Join),
			filesystem.WithPerm(filesystem.RwRR))

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})
}

func TestCopyDir(t *testing.T) {
	t.Run("error_no_dir", func(t *testing.T) {
		// Arrange
		srcdir := filepath.Join(os.TempDir(), "invalid")

		// Act
		err := filesystem.CopyDir(srcdir, t.TempDir())

		// Assert
		assert.ErrorContains(t, err, "failed to read directory")
	})

	t.Run("error_no_destdir", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		destdir := filepath.Join(os.TempDir(), "invalid", "dir")

		// Act
		err = filesystem.CopyDir(srcdir, destdir)

		// Assert
		assert.ErrorContains(t, err, "failed to create folder")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		srcsubdir := filepath.Join(srcdir, "sub", "dir")
		require.NoError(t, os.MkdirAll(srcsubdir, filesystem.RwxRxRxRx))
		srcsub := filepath.Join(srcsubdir, "file.txt")
		file, err = os.Create(srcsub)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		destdir := filepath.Join(os.TempDir(), "dir_test")
		t.Cleanup(func() {
			require.NoError(t, os.RemoveAll(destdir))
		})

		// Act
		err = filesystem.CopyDir(srcdir, destdir)

		// Assert
		assert.NoError(t, err)
		tests.AssertEqualDir(t, srcdir, destdir)
	})
}

func TestExists(t *testing.T) {
	t.Run("false_not_exists", func(t *testing.T) {
		// Arrange
		invalid := filepath.Join(os.TempDir(), "invalid")

		// Act
		exists := filesystem.Exists(invalid)

		// Assert
		assert.False(t, exists)
	})

	t.Run("true_exists", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		exists := filesystem.Exists(src)

		// Assert
		assert.True(t, exists)
	})
}
