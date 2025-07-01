package semver

import (
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
	parsing := 0
	major := make([]byte, 0)
	minor := make([]byte, 0)
	patch := make([]byte, 0)
	pre := make([]byte, 0)
	build := make([]byte, 0)

	for i := 0; i < len(input); i += 1 {
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
			case parsingMinor:
				minor = append(minor, b)
			case parsingPatch:
				patch = append(patch, b)
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

	majorInt, _ := strconv.Atoi(string(major))
	minorInt, _ := strconv.Atoi(string(minor))
	patchInt, _ := strconv.Atoi(string(patch))

	return &Version{
		major: majorInt,
		minor: minorInt,
		patch: patchInt,
		pre:   string(pre),
		build: strings.TrimSpace(string(build)),
	}, nil
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
