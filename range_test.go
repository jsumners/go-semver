package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange_RangeFromString(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title:    "simple empty string",
			input:    "",
			expected: ">=0.0.0",
		},
		{
			title:    "simple no operator",
			input:    "1.2.3",
			expected: "=1.2.3",
		},
		{
			title:    "simple two comparators",
			input:    ">1.0.0 <2.0.0",
			expected: ">1.0.0 <2.0.0",
		},
		{
			title:    "two sets present",
			input:    ">1.0.0 || >3 <=3.1.0",
			expected: ">1.0.0 || >3.0.0 <=3.1.0",
		},
		{
			// This case is made up nonsense. But it has shown to be interesting
			// while adding support for X-Ranges. The versions qualify as an x-range,
			// but each is prefixed with an operator that is contrary to what would
			// be implied by the x-range. The original thinking is that the input
			// would translate to `>1.0.0 || >2.0.0 || <5.0.0`. Since this is a weird
			// case, and likely not valid, we have accepted the wacky expansion as
			// valid and moved on. But this might be a good spot for improvement.
			title:    "three sets present",
			input:    ">1 || >2 || <5",
			expected: ">1.0.0 <2.0.0 || >2.0.0 <3.0.0 || <5.0.0 <6.0.0",
		},

		// Hyphen ranges
		{
			title:    "hyphen: basic inclusive set",
			input:    "1.2.3 - 2.3.4",
			expected: ">=1.2.3 <=2.3.4",
		},
		{
			title:    "hyphen: partial first (major & minor)",
			input:    "1.2 - 2.3.4",
			expected: ">=1.2.0 <=2.3.4",
		},
		{
			title:    "hyphen: partial first (major only)",
			input:    "1 - 2.3.4",
			expected: ">=1.0.0 <=2.3.4",
		},
		{
			title:    "hyphen: partial second (major & minor)",
			input:    "1.2.3 - 2.3",
			expected: ">=1.2.3 <2.4.0",
		},
		{
			title:    "hyphen: partial second (major only)",
			input:    "1.2.3 - 2",
			expected: ">=1.2.3 <3.0.0",
		},
		{
			title:    "hyphen: both partial (minor only)",
			input:    "1 - 2",
			expected: ">=1.0.0 <3.0.0",
		},

		// X ranges
		{
			title:    "x-range: any version *",
			input:    " * ", // spaces are intentional to verify they do not matter
			expected: ">=0.0.0",
		},
		{
			title:    "x-range: any version empty string",
			input:    "",
			expected: ">=0.0.0",
		},
		{
			title:    "x-range: major partial",
			input:    "1.x",
			expected: ">=1.0.0 <2.0.0",
		},
		{
			title:    "x-range: minor partial",
			input:    "1.2.x",
			expected: ">=1.2.0 <1.3.0",
		},
		{
			title:    "x-range: major only",
			input:    "1",
			expected: ">=1.0.0 <2.0.0",
		},

		// Tilde ranges:
		{
			title:    "tilde range: major, minor, patch",
			input:    "~1.2.3",
			expected: ">=1.2.3 <1.3.0",
		},
		{
			title:    "tilde range: major and minor",
			input:    "~1.2",
			expected: ">=1.2.0 <1.3.0",
		},
		{
			title:    "tilde range: major",
			input:    "~1",
			expected: ">=1.0.0 <2.0.0",
		},
		{
			title:    "tilde range: major 0, minor, patch",
			input:    "~0.2.3",
			expected: ">=0.2.3 <0.3.0",
		},
		{
			title:    "tilde range: major 0, minor",
			input:    "~0.2",
			expected: ">=0.2.0 <0.3.0",
		},
		{
			title:    "tilde range: major 0",
			input:    "~0",
			expected: ">=0.0.0 <1.0.0",
		},
		{
			title:    "tilde range: with pre",
			input:    "~1.2.3-beta.2",
			expected: ">=1.2.3-beta.2 <1.3.0",
		},

		// Caret ranges:
		{
			title:    "caret range: major, minor, patch",
			input:    "^1.2.3",
			expected: ">=1.2.3 <2.0.0",
		},
		{
			title:    "caret range: major 0, minor, patch",
			input:    "^0.2.3",
			expected: ">=0.2.3 <0.3.0",
		},
		{
			title:    "caret range: major 0, minor 0, patch",
			input:    "^0.0.3",
			expected: ">=0.0.3 <0.0.4",
		},
		{
			title:    "caret range: major, minor, patch, pre-release",
			input:    "^1.2.3-beta.2",
			expected: ">=1.2.3-beta.2 <2.0.0",
		},
		{
			title:    "caret range: major 0, minor 0, patch, pre-release",
			input:    "^0.0.3-beta",
			expected: ">=0.0.3-beta <0.0.4",
		},
		{
			title:    "caret range: major, minor, patch-x",
			input:    "^1.2.x",
			expected: ">=1.2.0 <2.0.0",
		},
		{
			title:    "caret range: major 0, minor 0, patch-x",
			input:    "^0.0.x",
			expected: ">=0.0.0 <0.1.0",
		},
		{
			title:    "caret range: major 0, minor 0",
			input:    "^0.0",
			expected: ">=0.0.0 <0.1.0",
		},
		{
			title:    "caret range: major, minor-x",
			input:    "^1.x",
			expected: ">=1.0.0 <2.0.0",
		},
		{
			title:    "caret range: major 0, minor-x",
			input:    "^0.x",
			expected: ">=0.0.0 <1.0.0",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			rng, err := RangeFromString(testCase.input)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, rng.String())
		})
	}

	t.Run("alpha char in range string", func(t *testing.T) {
		input := ">=A.0.1"
		res, err := RangeFromString(input)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, "encountered alpha character in range string: `>=A.0.1`", err.Error())
		assert.ErrorIs(t, err, ErrRangeAlpha)
	})
}
