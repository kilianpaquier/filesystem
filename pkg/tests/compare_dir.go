package tests

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

// readDirInMap reads a given input directory (and its subdirectories) and returns a map with filenames as keys and content (string) as values.
//
// Collision will occur in case a two files with the same name exists (between root and subdirectory).
func readDirInMap(srcdir string) (map[string]string, error) {
	files := map[string]string{}

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
		files[entry.Name()] = string(bytes)
	}
	return files, errors.Join(errs...)
}

// EqualDirOpt represents an option to be given to AssertEqualDir to tune its behavior.
type EqualDirOpt func(e *equalDirClient)

// IgnoreFunc represents a function to be given with WithIgnoreDiff to ignore a specific difference in AssertEqualDir.
type IgnoreFunc func(filename string, item diffmatchpatch.Diff) bool

// equalDirClient represents the client option of AssertEqualDir.
type equalDirClient struct {
	ignore []IgnoreFunc
}

// WithIgnoreDiff is an option to give to AssertEqualDir to ignore specific differences during execution.
func WithIgnoreDiff(ignore IgnoreFunc) EqualDirOpt {
	return func(e *equalDirClient) {
		e.ignore = append(e.ignore, ignore)
	}
}

// AssertEqualDir compares expected an actual directories (and their subdirectories).
//
// It will fail with t in case a file is missing in actual, a file is present in actual but not in expected and if the content of any file in actual is not the same as its peer in expected.
func AssertEqualDir(t testing.TB, expected, actual string, opts ...EqualDirOpt) {
	// read all files in expected directory
	expectedFiles, err := readDirInMap(expected)
	assert.NoError(t, err, "failed to completely read expected %s folder and its children", expected)

	// read all files in actual directory
	actualFiles, err := readDirInMap(actual)
	assert.NoError(t, err, "failed to completely read actual %s folder and its children", actual)

	client := &equalDirClient{}
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	// check all expected contents against actual contents
	for expectedFilename, expectedContent := range expectedFiles {
		actualContent, ok := actualFiles[expectedFilename]
		assert.True(t, ok, "%s missing from actual directory", expectedFilename)
		if !ok {
			continue
		}

		matcher := diffmatchpatch.New()
		diffs := matcher.DiffMain(expectedContent, actualContent, false)

		// filter retrieved diffs
		diffs = client.filterDiffs(expectedFilename, diffs)
		assert.Len(t, diffs, 0, "%s is different between expected and actual folder: %s", expectedFilename, matcher.DiffPrettyText(diffs))
	}

	for actualFilename := range actualFiles {
		_, ok := expectedFiles[actualFilename]
		assert.True(t, ok, "%s is present in actual directory but not in expected one", actualFilename)
	}
}

// filterDiffs filters input diffs with the equalDirClient options.
func (e equalDirClient) filterDiffs(filename string, diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	filtered := make([]diffmatchpatch.Diff, 0, len(diffs))
	for _, diff := range diffs {
		// ignore diff equals
		// ignore windows / linux diffs
		if diff.Type == diffmatchpatch.DiffEqual || diff.Text == "\r" {
			continue
		}

		// ignore some diffs according to input options
		isIgnored := false
		for _, ignore := range e.ignore {
			if ignore(filename, diff) {
				isIgnored = true
				break
			}
		}

		// add diff to filtered diffs in case it has to be kept
		if !isIgnored {
			filtered = append(filtered, diff)
		}
	}
	return filtered
}
