package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FSOption represents a function taking an opt client to use filesysem package functions.
type FSOption func(fsOpt *fsOpt)

// WithFS specifies a FS to read files instead of os filesystem.
func WithFS(fsys FS) FSOption {
	return func(fsOpt *fsOpt) {
		fsOpt.fsys = fsys
	}
}

// Join represents a function to join multiple elements between them.
type Join func(elems ...string) string

// WithJoin specifies a specific function to join a srcdir with one of its files in CopyDir.
func WithJoin(join Join) FSOption {
	return func(fsOpt *fsOpt) {
		fsOpt.join = join
	}
}

type fsOpt struct {
	fsys FS
	join Join
}

func newFSOpt(opts ...FSOption) *fsOpt {
	o := &fsOpt{}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	if o.fsys == nil {
		o.fsys = OS()
	}
	if o.join == nil {
		o.join = filepath.Join
	}
	return o
}

// CopyFile copies a provided file from src to dest with a default permission of 0o644. It fails if it's a directory.
func CopyFile(src, dest string, opts ...FSOption) error {
	return CopyFileWithPerm(src, dest, RwRR, opts...)
}

// CopyDir copies recursively a provided directory as destdir. It fails if it's a file.
func CopyDir(srcdir, destdir string, opts ...FSOption) error {
	o := newFSOpt(opts...)

	if err := os.Mkdir(destdir, RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create folder %s: %w", destdir, err)
	}

	entries, err := o.fsys.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		src := o.join(srcdir, entry.Name())
		dest := filepath.Join(destdir, entry.Name())

		// handle directories
		if entry.IsDir() {
			errs = append(errs, CopyDir(src, dest, opts...))
			continue
		}

		// handle files
		errs = append(errs, CopyFile(src, dest, opts...))
	}
	return errors.Join(errs...)
}

// CopyFileWithPerm copies a provided file from src to dest with specific permissions. It fails if it's a directory.
func CopyFileWithPerm(src, dest string, perm fs.FileMode, opts ...FSOption) error {
	o := newFSOpt(opts...)

	// read file from fsys (OperatingFS or specific fsys)
	bytes, err := o.fsys.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", src, err)
	}

	// write file destination into OperatingFS
	if err := os.WriteFile(dest, bytes, perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", dest, err)
	}
	return nil
}

// Exists returns a boolean indicating whether the provided input src exists or not.
func Exists(src string, opts ...FSOption) bool {
	o := newFSOpt(opts...)

	// read file from fsys (OperatingFS or specific fsys)
	file, err := o.fsys.Open(src)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}
