package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion_VersionFromString(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected *Version
	}{
		{
			// `*` => `0.0.0`
			title:    "any version",
			input:    " * ",
			expected: &Version{partial: true},
		},
		{
			title: "parses version with pre-release and build",
			input: "1.2.3-alpha.1+build.2  ", // Spaces are on purpose.
			expected: &Version{
				major:       1,
				minor:       2,
				patch:       3,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "alpha.1",
				build:       "build.2",
			},
		},
		{
			title:    "parses simple major",
			input:    "42",
			expected: &Version{major: 42, majorParsed: true, partial: true},
		},
		{
			title: "parses major.minor",
			input: "42.24",
			expected: &Version{
				major:       42,
				minor:       24,
				majorParsed: true,
				minorParsed: true,
				partial:     true,
			},
		},
		{
			title: "parses a major.minor.patch",
			input: "42.24.9",
			expected: &Version{
				major:       42,
				minor:       24,
				patch:       9,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
			},
		},
		{
			title:    "parses simple major version with pre",
			input:    "1-pre.1",
			expected: &Version{major: 1, majorParsed: true, pre: "pre.1", partial: true},
		},

		// Pre-release examples from spec:
		{
			title: "pre-release 1.0.0-alpha",
			input: "1.0.0-alpha",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "alpha",
			},
		},
		{
			title: "pre-release 1.0.0-alpha.1",
			input: "1.0.0-alpha.1",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "alpha.1",
			},
		},
		{
			title: "pre-release 1.0.0-0.3.7",
			input: "1.0.0-0.3.7",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "0.3.7",
			},
		},
		{
			title: "pre-release 1.0.0-x.7.z.92",
			input: "1.0.0-x.7.z.92",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "x.7.z.92",
			},
		},
		{
			title: "pre-release 1.0.0-x-y-z.--",
			input: "1.0.0-x-y-z.--",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "x-y-z.--",
			},
		},

		// Build examples from spec:
		{
			title: "build 1.0.0+alpha+001",
			input: "1.0.0+alpha+001",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				build:       "alpha+001",
			},
		},
		{
			title: "build 1.0.0+20130313144700",
			input: "1.0.0+20130313144700",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				build:       "20130313144700",
			},
		},
		{
			title: "build 1.0.0-beta+exp.sha.5114f85",
			input: "1.0.0-beta+exp.sha.5114f85",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				pre:         "beta",
				build:       "exp.sha.5114f85",
			},
		},
		{
			title: "build 1.0.0+21AF26D3----117B344092BD",
			input: "1.0.0+21AF26D3----117B344092BD",
			expected: &Version{
				major:       1,
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
				build:       "21AF26D3----117B344092BD",
			},
		},

		// x-range versions:
		{
			title: "x in patch position",
			input: "1.1.x",
			expected: &Version{
				major:       1,
				minor:       1,
				patch:       0,
				majorParsed: true,
				minorParsed: true,
				patchParsed: false,
				partial:     true,
			},
		},
		{
			title: "x in minor and patch position",
			input: "1.x",
			expected: &Version{
				major:       1,
				minor:       0,
				patch:       0,
				majorParsed: true,
				minorParsed: false,
				patchParsed: false,
				partial:     true,
			},
		},
		{
			title: "x in major position",
			input: "X",
			expected: &Version{
				major:       0,
				minor:       0,
				patch:       0,
				majorParsed: false,
				minorParsed: false,
				patchParsed: false,
				partial:     true,
			},
		},
		{
			title: "x in minor and patch position with pre",
			input: "1.x-abc",
			expected: &Version{
				major:       1,
				minor:       0,
				patch:       0,
				pre:         "abc",
				majorParsed: true,
				minorParsed: false,
				patchParsed: false,
				partial:     true,
			},
		},
		{
			title: "x in minor and patch position with build",
			input: "1.x+abc",
			expected: &Version{
				major:       1,
				minor:       0,
				patch:       0,
				build:       "abc",
				majorParsed: true,
				minorParsed: false,
				patchParsed: false,
				partial:     true,
			},
		},
		{
			title: "x in minor and patch position with pre and build",
			input: "1.x-abc+def",
			expected: &Version{
				major:       1,
				minor:       0,
				patch:       0,
				pre:         "abc",
				build:       "def",
				majorParsed: true,
				minorParsed: false,
				patchParsed: false,
				partial:     true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			ver, err := VersionFromString(testCase.input)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, ver)
		})
	}
}

func TestVersion_Satisfies(t *testing.T) {
	testCases := []struct {
		title       string
		expected    bool
		version     string
		targetRange string
	}{
		{
			title:       "equal (true)",
			expected:    true,
			version:     "1.0.0",
			targetRange: "=1.0.0",
		},
		{
			title:       "equal (false)",
			expected:    false,
			version:     "1.1.0",
			targetRange: "=1.0.0",
		},
		{
			title:       "less than (true)",
			expected:    true,
			version:     "0.9.0",
			targetRange: "<1.0.0",
		},
		{
			title:       "less than (false)",
			expected:    false,
			version:     "1.1.0",
			targetRange: "<1.0.0",
		},
		{
			title:       "less than equal (true)",
			expected:    true,
			version:     "1.0.0",
			targetRange: "<=1.0.0",
		},
		{
			title:       "less than equal (false)",
			expected:    false,
			version:     "1.1.0",
			targetRange: "<=1.0.0",
		},
		{
			title:       "greater than (true)",
			expected:    true,
			version:     "2.0.0",
			targetRange: ">1.0.0",
		},
		{
			title:       "greater than (false)",
			expected:    false,
			version:     "0.9.0",
			targetRange: ">1.0.0",
		},
		{
			title:       "greater than equal (true)",
			expected:    true,
			version:     "1.0.0",
			targetRange: ">=1.0.0",
		},
		{
			title:       "greater than equal (false)",
			expected:    false,
			version:     "0.9.0",
			targetRange: ">=1.0.0",
		},

		{
			title:       "within lower and upper (true)",
			expected:    true,
			version:     "1.5.0",
			targetRange: ">1 <2",
		},
		{
			title:       "lower, outside lower and upper",
			expected:    false,
			version:     "0.0.1",
			targetRange: ">=1.2 <2",
		},
		{
			title:       "upper, outside lower and upper",
			expected:    false,
			version:     "2.1",
			targetRange: ">=1.2 <2",
		},

		{
			title:       "or-ed: in range, lower",
			expected:    true,
			version:     "0.4.0",
			targetRange: "<0.5 || >0.8",
		},
		{
			title:       "or-ed: in range, upper",
			expected:    true,
			version:     "0.8.9",
			targetRange: "<0.5 || >0.8",
		},
		{
			title:       "or-ed: in range, middle",
			expected:    true,
			version:     "0.6.0",
			targetRange: ">0.5 || <0.8",
		},
		{
			title:       "or-ed: in range, third satisfies",
			expected:    true,
			version:     "3",
			targetRange: "=1.0.0 || =2.0.0 || =3.0.0",
		},
		{
			// `>0.8` expands to `>0.8.0 <0.9.0`
			title:       "or-ed: not in range, upper",
			expected:    false,
			version:     "0.9.0",
			targetRange: "<0.5 || >0.8",
		},
		{
			title:       "or-ed: none satisfies",
			expected:    false,
			version:     "4",
			targetRange: "=1.0.0 || =2.0.0 || =3.0.0",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			v, _ := VersionFromString(testCase.version)
			r, _ := RangeFromString(testCase.targetRange)
			found := v.Satisfies(r)
			assert.Equal(t, testCase.expected, found)
		})
	}
}

func TestVersion_Equals(t *testing.T) {
	v1, _ := VersionFromString("1")
	v2, _ := VersionFromString("1")
	assert.Equal(t, true, v1.Equals(v2))

	v2, _ = VersionFromString("2")
	assert.Equal(t, false, v1.Equals(v2))
}

func TestVersion_Greater(t *testing.T) {
	v1, _ := VersionFromString("1")
	v2, _ := VersionFromString("1")
	assert.Equal(t, false, v1.Greater(v2))

	v2, _ = VersionFromString("0")
	assert.Equal(t, true, v1.Greater(v2))
}

func TestVersion_GreaterThanEquals(t *testing.T) {
	v1, _ := VersionFromString("1")
	v2, _ := VersionFromString("1")
	assert.Equal(t, true, v1.GreaterThanEquals(v2))

	v2, _ = VersionFromString("2")
	assert.Equal(t, false, v1.GreaterThanEquals(v2))
}

func TestVersion_Less(t *testing.T) {
	v1, _ := VersionFromString("1")
	v2, _ := VersionFromString("1")
	assert.Equal(t, false, v1.Less(v2))

	v2, _ = VersionFromString("0")
	assert.Equal(t, true, v2.Less(v1))
}

func TestVersion_LessThanEquals(t *testing.T) {
	v1, _ := VersionFromString("1")
	v2, _ := VersionFromString("1")
	assert.Equal(t, true, v1.LessThanEquals(v2))

	v2, _ = VersionFromString("2")
	assert.Equal(t, true, v1.LessThanEquals(v2))
}
