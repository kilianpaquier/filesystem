package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

type osFS struct{}

var operating FS = &osFS{}

// OS returns an implementation of FS for the current filesystem.
func OS() FS {
	return operating
}

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func (*osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func (*osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func (*osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Join joins any number of path elements into a single path,
// separating them with an OS specific [Separator]. Empty elements
// are ignored. The result is Cleaned. However, if the argument
// list is empty or all its elements are empty, Join returns
// an empty string.
// On Windows, the result will only be a UNC path if the first
// non-empty element is a UNC path.
func (*osFS) Join(elems ...string) string {
	return filepath.Join(elems...)
}
