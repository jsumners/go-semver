package semver

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrVersionParseFailure = errors.New("failed to parse version string")

type Version struct {
	major int
	minor int
	patch int

	/* We need a way to differentiate the basic zero value of components from
	a zero value read from a provided version string.
	*/
	// majorParsed indicates the value of major was parsed from the version string.
	majorParsed bool
	// minorParsed indicates the value of minor was parsed from the version string.
	minorParsed bool
	// patchParsed indicates the value of patch was parsed from the version string.
	patchParsed bool

	// partial indicates that the version had an "x" or "X" in one of the three
	// primary positions.
	partial bool

	pre   string
	build string
}

const (
	parsingMajor = iota
	parsingMinor
	parsingPatch
	parsingPre
	parsingBuild
)

func VersionFromString(input string) (*Version, error) {
	return VersionFromBytes([]byte(input))
}

func VersionFromBytes(input []byte) (*Version, error) {
	if bytes.Equal(bytes.TrimSpace(input), []byte("*")) == true {
		return &Version{partial: true}, nil
	}

	parsing := 0
	major := make([]byte, 0)
	minor := make([]byte, 0)
	patch := make([]byte, 0)
	pre := make([]byte, 0)
	build := make([]byte, 0)

	isPartial := true
	majorParsed := false
	minorParsed := false
	patchParsed := false

	strLen := len(input)
	for i := 0; i < strLen; i += 1 {
		b := input[i] // current byte
		c := int(b)   // current byte as a "character" because we can't compare bytes

		if i == 0 && (c == capitalV || c == lowerV) {
			// We don't care about a leading `V` or `v`.
			continue
		}

		if parsing < parsingPre && (isIntegerChar(b) == true) {
			switch parsing {
			case parsingMajor:
				major = append(major, b)
				majorParsed = true
			case parsingMinor:
				minor = append(minor, b)
				minorParsed = true
			case parsingPatch:
				patch = append(patch, b)
				patchParsed = true
				isPartial = false
			}
			continue
		} else if parsing < parsingPre && (isXRangeChar(b) == true) {
			switch parsing {
			case parsingMajor:
				major = append(major, byte(0))
				parsing += 3
			case parsingMinor:
				minor = append(minor, byte(0))
				parsing += 2
			case parsingPatch:
				patch = append(minor, byte(0))
				parsing += 1
			}

			// We need to inspect the next character, if it exits, in order to
			// prepare our state for the next iteration.
			if i+1 == strLen {
				break
			} else if int(input[i+1]) == dash {
				// We need to advance the position or else the dash separator will
				// be included in the pre-release string.
				i += 1
			} else if int(input[i+1]) == plus {
				// We don't advance the position here because our test for build
				// strings hinges on both the parsing indicator and the presence
				// of the plus (`+`) character.
				parsing += 1
			}

			continue
		}

		// Increment parsing as we encounter periods (`.`).
		if parsing < parsingPre && c == dot {
			parsing += 1
			continue
		}

		// Build strings are the final component.
		if parsing <= parsingBuild && c == plus {
			// todo: we probably want to continue with byte-wise parsing in order
			// to validate that the build string is comprised of valid characters
			build = input[i+1:]
			break
		}

		// We need "< parsingPre" because a dash (`-`) is a valid character
		// in a pre-release identifier. Rather frustrating that it is the separator
		// between the patch number and the pre-release identifier as well as a
		// valid identifier character, but that is the spec.
		if parsing < parsingPre && c == dash {
			parsing = parsingPre
			continue
		}

		if c == dash ||
			c == dot || // since we are just stringifying any pre or build tags, include dot
			(isIntegerChar(b) == true) || // [0-9]
			(isCapitalChar(b) == true) || // [A-Z]
			(isLowerChar(b) == true) { // [a-z]
			if parsing == parsingPre {
				pre = append(pre, b)
				continue
			}
		}
	}

	version := &Version{
		pre:     string(pre),
		build:   strings.TrimSpace(string(build)),
		partial: isPartial,
	}
	if len(major) > 0 {
		majorInt, _ := strconv.Atoi(string(major))
		version.major = majorInt
		version.majorParsed = majorParsed
	} else {
		version.major = 0
	}
	if len(minor) > 0 {
		minorInt, _ := strconv.Atoi(string(minor))
		version.minor = minorInt
		version.minorParsed = minorParsed
	} else {
		version.minor = 0
	}
	if len(patch) > 0 {
		patchInt, _ := strconv.Atoi(string(patch))
		version.patch = patchInt
		version.patchParsed = patchParsed
	} else {
		version.patch = 0
	}

	return version, nil
}

func (v *Version) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch))
	if v.pre != "" {
		builder.WriteString(fmt.Sprintf("-%s", v.pre))
	}
	if v.build != "" {
		builder.WriteString(fmt.Sprintf("+%s", v.build))
	}
	return builder.String()
}

// Compare is a convenience method for the generic [Compare] function.
func (v *Version) Compare(ver *Version) int {
	return Compare(v, ver)
}

// Equals is a convenience method to determine if the version is
// equivalent to another [Version] instance. This is a by value comparison.
func (v *Version) Equals(ver *Version) bool {
	return v.Compare(ver) == 0
}

// Greater is a convenience method to determine if the version is
// a higher version than the provided [Version].
func (v *Version) Greater(ver *Version) bool {
	return v.Compare(ver) == 1
}

// GreaterThanEquals is a convenience method to determine if the version is
// equal to or greater than the provided [Version].
func (v *Version) GreaterThanEquals(ver *Version) bool {
	return v.Compare(ver) >= 0
}

// Less is a convenience method to determine if the version is
// a lower version than the provided [Version].
func (v *Version) Less(ver *Version) bool {
	return v.Compare(ver) == -1
}

// LessThanEquals is a convenience method to determine if the version is
// equal to or less than the provided [Version].
func (v *Version) LessThanEquals(ver *Version) bool {
	return v.Compare(ver) <= 0
}

// Satisfies determines if the version is covered by the provided
// [Range].
func (v *Version) Satisfies(r *Range) bool {
	var result bool

	for _, set := range r.comparators {
		if result == true {
			// We've iterated through at least one comparator set and the outcome
			// was that the version satisfies the range defined by that set. According
			// to the spec, this means the condition has been satisfied regardless
			// of what any other set in the range would indicate.
			break
		}

		if set.one != nil && set.two != nil {
			a := inRange(v, set.one)
			b := inRange(v, set.two)
			result = a == true && b == true
			continue
		}

		if set.one != nil {
			a := inRange(v, set.one)
			if a == true {
				result = a
			}
		}
		if set.two != nil {
			b := inRange(v, set.two)
			if b == true {
				result = b
			}
		}
	}

	return result
}

func inRange(ver *Version, comp *Comparator) bool {
	compareResult := ver.Compare(comp.version)
	switch comp.operator {
	case OperatorEqual:
		return compareResult == 0
	case OperatorLessThan:
		return compareResult == -1
	case OperatorLessThanEqual:
		return compareResult == 0 || compareResult == -1
	case OperatorGreaterThan:
		return compareResult == 1
	case OperatorGreaterThanEqual:
		return compareResult == 1 || compareResult == 0
	default:
		return false
	}
}
