package semver

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var ErrRangeAlpha = errors.New("encountered alpha character in range string")

type RangeOperator int

const (
	OperatorEqual RangeOperator = iota
	OperatorLessThan
	OperatorLessThanEqual
	OperatorGreaterThan
	OperatorGreaterThanEqual
	OperatorUnknown
)

// TODO: support "advanced range syntax" (e.g. ^1.0.0)
func (r RangeOperator) String() string {
	switch r {
	case OperatorEqual:
		return "="
	case OperatorLessThan:
		return "<"
	case OperatorLessThanEqual:
		return "<="
	case OperatorGreaterThan:
		return ">"
	case OperatorGreaterThanEqual:
		return ">="
	default:
		return "<>"
	}
}

func RangeOperatorFromBytes(input []byte) RangeOperator {
	switch {
	case bytes.Equal(input, []byte("=")):
		return OperatorEqual
	case bytes.Equal(input, []byte("<")):
		return OperatorLessThan
	case bytes.Equal(input, []byte("<=")):
		return OperatorLessThanEqual
	case bytes.Equal(input, []byte(">")):
		return OperatorGreaterThan
	case bytes.Equal(input, []byte(">=")):
		return OperatorGreaterThanEqual
	default:
		return OperatorUnknown
	}
}

type Comparator struct {
	operator      RangeOperator
	version       *Version
	operatorBytes []byte
	versionBytes  []byte
}

func newComparator() *Comparator {
	return &Comparator{
		operatorBytes: make([]byte, 0),
		versionBytes:  make([]byte, 0),
	}
}

func (c *Comparator) String() string {
	builder := strings.Builder{}
	builder.WriteString(c.operator.String())
	builder.WriteString(c.version.String())
	return builder.String()
}

type ComparatorSet struct {
	one *Comparator
	two *Comparator
}

type Range struct {
	comparators []ComparatorSet
}

func RangeFromString(input string) (*Range, error) {
	return RangeFromBytes([]byte(input))
}

func RangeFromBytes(input []byte) (*Range, error) {
	if bytes.Equal(input, []byte("")) == true {
		// An empty string is a special case range that maps
		// to ">=0.0.0".
		c := &Comparator{
			operator: OperatorGreaterThanEqual,
			version: &Version{
				majorParsed: true,
				minorParsed: true,
				patchParsed: true,
			},
			operatorBytes: nil,
			versionBytes:  nil,
		}
		set := ComparatorSet{one: c}
		return &Range{comparators: []ComparatorSet{set}}, nil
	}

	setsToParse := bytes.Split(input, []byte("||"))

	comparators := make([]ComparatorSet, 0)
	for _, set := range setsToParse {
		set = bytes.TrimSpace(set)

		range1, range2, found := bytes.Cut(set, []byte(" - "))
		if found == true {
			// We have a hyphen range, e.g. `1.0.0 - 2.0.0`.
			comparatorSet, err := parseHyphenRange(range1, range2)
			if err != nil {
				return nil, err
			}
			comparators = append(comparators, comparatorSet)
			continue
		}

		range1, range2, found = bytes.Cut(set, []byte(" "))
		if found == true {
			// We have a basic range separated by a space, e.g. `1.0.0 2.0.0`.
			comparatorSet, err := parseBasicRange(range1, range2)
			if err != nil {
				return nil, err
			}
			comparators = append(comparators, comparatorSet)
			continue
		}

		// We have a simple range, e.g. `=1.0.0`.
		comparator, err := parseComparator(set)
		if err != nil {
			return nil, err
		}
		comparators = append(comparators, ComparatorSet{one: comparator})
	}

	return &Range{comparators: comparators}, nil
}

func parseHyphenRange(r1 []byte, r2 []byte) (ComparatorSet, error) {
	c1, err := parseComparator(bytes.TrimSpace(r1))
	if err != nil {
		return ComparatorSet{}, nil
	}

	c2, err := parseComparator(bytes.TrimSpace(r2))
	if err != nil {
		return ComparatorSet{}, nil
	}

	c1.operator = OperatorGreaterThanEqual
	if c2.version.patchParsed == true {
		c2.operator = OperatorLessThanEqual
	} else if c2.version.minorParsed == true {
		c2.version.minor += 1
		c2.operator = OperatorLessThan
	} else {
		c2.version.major += 1
		c2.operator = OperatorLessThan
	}

	return ComparatorSet{c1, c2}, nil
}

func parseBasicRange(r1 []byte, r2 []byte) (ComparatorSet, error) {
	c1, err := parseComparator(bytes.TrimSpace(r1))
	if err != nil {
		return ComparatorSet{}, err
	}

	c2, err := parseComparator(bytes.TrimSpace(r2))
	if err != nil {
		return ComparatorSet{}, err
	}

	return ComparatorSet{c1, c2}, nil
}

func parseComparator(r []byte) (*Comparator, error) {
	comparator := newComparator()

	for i := 0; i < len(r); i += 1 {
		b := r[i]
		if isAlphaChar(b) == true {
			// TODO: handle `x`, `X`, and `*`
			return nil, fmt.Errorf("%w: `%s`", ErrRangeAlpha, r)
		}

		if isOperatorChar(b) {
			comparator.operatorBytes = append(comparator.operatorBytes, b)
			continue
		}

		comparator.versionBytes = r[i:]
		break
	}

	err := finalizeComparator(comparator)
	if err != nil {
		return nil, err
	}

	return comparator, nil
}

func finalizeComparator(c *Comparator) error {
	if len(c.operatorBytes) > 0 {
		c.operator = RangeOperatorFromBytes(c.operatorBytes)
		c.operatorBytes = make([]byte, 0)
	}
	if len(c.versionBytes) > 0 {
		ver, err := VersionFromBytes(c.versionBytes)
		if err != nil {
			return fmt.Errorf("%w: `%s`", err, c.versionBytes)
		}
		c.version = ver
		c.versionBytes = make([]byte, 0)
	}
	return nil
}

func (r *Range) String() string {
	builder := strings.Builder{}
	for _, set := range r.comparators {
		builder.WriteString(set.one.String())
		if set.two != nil {
			builder.WriteString(fmt.Sprintf(" %s", set.two))
		}
		builder.WriteString(" || ")
	}
	s := builder.String()
	return strings.TrimSuffix(s, " || ")
}
