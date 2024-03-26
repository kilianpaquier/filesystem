package tests

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diffs returns the slice of diffs between the two input contents
// and a boolean indicating whether they are identical or not.
func Diffs(expected, actual string) ([]diffmatchpatch.Diff, bool) {
	matcher := diffmatchpatch.New()

	// compute diffs between expected content and actual content with cleanup runs
	diffs := matcher.DiffMain(expected, actual, false)
	diffs = matcher.DiffCleanupSemantic(diffs)
	diffs = matcher.DiffCleanupEfficiency(diffs)

	// find the first not equal diff and not being the CRLF / LF diff between windows and linux
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual && diff.Text != "\r" {
			return diffs, false
		}
	}
	return diffs, true
}
