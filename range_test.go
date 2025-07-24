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
			title:    "three sets present",
			input:    ">1 || >2 || <5",
			expected: ">1.0.0 || >2.0.0 || <5.0.0",
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
