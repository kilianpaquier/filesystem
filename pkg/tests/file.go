package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FilterCarriage returns the input slice of bytes with \r character.
func FilterCarriage(bytes []byte) []byte {
	result := make([]byte, 0, len(bytes))
	for _, b := range bytes {
		if b != 13 {
			result = append(result, b)
		}
	}
	return result
}

// AssertEqualFile compares expected and actual files.
//
// It will fail with t if one of the file cannot be read or if their content is not identical.
func AssertEqualFile(t testing.TB, expected, actual string) {
	expectedBytes, err := os.ReadFile(expected)
	assert.NoError(t, err, "failed to read %s", expected)

	actualBytes, err := os.ReadFile(actual)
	assert.NoError(t, err, "failed to read %s", actual)

	diffs := Diff(expected, FilterCarriage(expectedBytes), actual, FilterCarriage(actualBytes))
	if len(diffs) > 0 {
		assert.Fail(t, actual+" is different from expected", string(diffs))
	}
}
