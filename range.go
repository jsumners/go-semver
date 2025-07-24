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

		parsingVersion := false
		comparatorSet := ComparatorSet{}
		comparator := newComparator()
		for i := 0; i < len(set); i += 1 {
			b := set[i]
			if isAlphaChar(b) == true {
				return nil, fmt.Errorf("%w: `%s`", ErrRangeAlpha, input)
			}

			if isOperatorChar(b) && parsingVersion == false {
				comparator.operatorBytes = append(comparator.operatorBytes, b)
				continue
			} else if isOperatorChar(b) && parsingVersion == true {
				// The first comparator has been parsed. For example, in the range
				// `>1.0.0 <=1.9.0`, at this point ">1.0.0" has been parsed, and we
				// need to start parsing "<=1.9.0".
				comparatorSet.one = comparator
				comparator = newComparator()
				comparator.operatorBytes = append(comparator.operatorBytes, b)
				parsingVersion = false
				continue
			}

			comparator.versionBytes = append(comparator.versionBytes, b)
			parsingVersion = true
		}

		if comparatorSet.one == nil {
			// We parsed through the whole set and never assigned a comparator.
			// Thus, it must have been a single comparator, e.g. `>=1.2.3`.
			// Therefore, we must assign it.
			comparatorSet.one = comparator
		} else if comparatorSet.two == nil {
			comparatorSet.two = comparator
		}

		if comparatorSet.one != nil {
			err := finalizeComparator(comparatorSet.one)
			if err != nil {
				return nil, err
			}
		}
		if comparatorSet.two != nil {
			err := finalizeComparator(comparatorSet.two)
			if err != nil {
				return nil, err
			}
		}

		comparators = append(comparators, comparatorSet)
	}

	return &Range{comparators: comparators}, nil
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
