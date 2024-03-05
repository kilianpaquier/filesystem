package filesystem

import (
	"io/fs"
)

// FS represents a filesystem with required minimal functions like Open, ReadDir and ReadFile.
type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}
