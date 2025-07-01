package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Compare(t *testing.T) {
	testCases := []struct {
		a        Version
		b        Version
		expected int
	}{
		{a: Version{major: 2}, b: Version{major: 1}, expected: 1},
		{a: Version{major: 1}, b: Version{major: 2}, expected: -1},
		{a: Version{major: 1}, b: Version{major: 1}, expected: 0},

		{a: Version{major: 1, minor: 2}, b: Version{major: 1, minor: 1}, expected: 1},
		{a: Version{major: 1, minor: 1}, b: Version{major: 1, minor: 2}, expected: -1},
		{a: Version{major: 1, minor: 1}, b: Version{major: 1, minor: 1}, expected: 0},

		{a: Version{major: 1, minor: 1, patch: 2}, b: Version{major: 1, minor: 1, patch: 1}, expected: 1},
		{a: Version{major: 1, minor: 1, patch: 1}, b: Version{major: 1, minor: 1, patch: 2}, expected: -1},
		{a: Version{major: 1, minor: 1, patch: 1}, b: Version{major: 1, minor: 1, patch: 1}, expected: 0},
	}

	for _, testCase := range testCases {
		found := Compare(&testCase.a, &testCase.b)
		assert.Equal(t, testCase.expected, found)
	}
}
