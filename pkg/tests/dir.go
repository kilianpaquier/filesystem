package tests

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// readDirInMap reads a given input directory (and its subdirectories) and returns a map with filenames as keys and content (string) as values.
//
// Collision will occur in case a two files with the same name exists (between root and subdirectory).
func readDirInMap(srcdir string) (map[string][]byte, error) {
	files := map[string][]byte{}

	entries, err := os.ReadDir(srcdir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", srcdir, err)
	}

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		src := filepath.Join(srcdir, entry.Name())

		// handle directories
		if entry.IsDir() {
			sub, err := readDirInMap(src)
			if err != nil {
				errs = append(errs, err) // only case of error is if reading an entry fails
			}

			for filename, content := range sub {
				// NOTE collision on identical filenames between root and subdirectories
				files[filename] = content
			}
			continue
		}

		// handle files
		bytes, err := os.ReadFile(src)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read %s: %w", src, err))
			continue
		}
		files[entry.Name()] = FilterCarriage(bytes)
	}
	return files, errors.Join(errs...)
}

// AssertEqualDir compares expected an actual directories (and their subdirectories).
//
// It will fail with t in case a file is missing in actual,
// a file is present in actual but not in expected
// or if the content of any file in actual is not the same as its peer in expected.
func AssertEqualDir(t testing.TB, expected, actual string) {
	// read all files in expected directory
	expectedFiles, err := readDirInMap(expected)
	assert.NoError(t, err, "failed to completely read expected %s folder and its children", expected)

	// read all files in actual directory
	actualFiles, err := readDirInMap(actual)
	assert.NoError(t, err, "failed to completely read actual %s folder and its children", actual)

	// check all expected contents against actual contents
	for filename, expectedBytes := range expectedFiles {
		actualBytes, ok := actualFiles[filename]
		assert.True(t, ok, "%s missing from actual directory", filename)
		if !ok {
			continue
		}

		diffs := Diff(filename, expectedBytes, filename, actualBytes)
		if len(diffs) > 0 {
			assert.Fail(t, filename+" is different from expected", string(diffs))
		}
	}

	// check that there're no actual files that aren't present in expected files
	for filename := range actualFiles {
		_, ok := expectedFiles[filename]
		assert.True(t, ok, "%s is present in actual directory but not in expected one", filename)
	}
}
