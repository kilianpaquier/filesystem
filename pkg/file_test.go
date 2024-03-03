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
	src := filepath.Join(tmp, "filename.txt")
	dest := filepath.Join(tmp, "filename-copy.txt")

	err := os.WriteFile(src, []byte("hey file"), filesystem.RwRR)
	require.NoError(t, err)

	t.Run("error_src_not_exists", func(t *testing.T) {
		// Arrange
		src := filepath.Join(tmp, "invalid.txt")

		// Act
		err := filesystem.CopyFile(src, dest)

		// Assert
		assert.ErrorContains(t, err, "failed to read "+src)
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
		err := filesystem.CopyFile(src, dest, filesystem.WithFS(filesystem.OS()))

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})
}

func TestCopyDir(t *testing.T) {
	t.Run("error_no_dir", func(t *testing.T) {
		// Arrange
		invalid := filepath.Join(os.TempDir(), "invalid")

		// Act
		err := filesystem.CopyDir(invalid, t.TempDir())

		// Assert
		assert.ErrorContains(t, err, "failed to read directory")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		dir := filepath.Join(srcdir, "path", "to", "dir")
		require.NoError(t, os.MkdirAll(dir, filesystem.RwxRxRxRx))

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

func TestCopyFileWithPerm(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "filename.txt")
	dest := filepath.Join(tmp, "filename-copy.txt")

	err := os.WriteFile(src, []byte("hey file"), filesystem.RwRR)
	require.NoError(t, err)

	t.Run("error_src_not_exists", func(t *testing.T) {
		// Arrange
		src := filepath.Join(tmp, "invalid.txt")

		// Act
		err := filesystem.CopyFileWithPerm(src, dest, filesystem.RwRwRw)

		// Assert
		assert.ErrorContains(t, err, "failed to read "+src)
		assert.NoFileExists(t, dest)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		err := filesystem.CopyFileWithPerm(src, dest, filesystem.RwRwRw)

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		err := filesystem.CopyFileWithPerm(src, dest, filesystem.RwRwRw, filesystem.WithFS(filesystem.OS()))

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
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
