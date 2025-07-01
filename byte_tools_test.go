package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isAlphaChar(t *testing.T) {
	input := "aB1"
	res := isAlphaChar(input[0])
	assert.Equal(t, true, res)

	res = isAlphaChar(input[1])
	assert.Equal(t, true, res)

	res = isAlphaChar(input[2])
	assert.Equal(t, false, res)
}

func Test_isOperatorChar(t *testing.T) {
	input := "<=<|"
	res := isOperatorChar(input[0])
	assert.Equal(t, true, res)

	res = isOperatorChar(input[1])
	assert.Equal(t, true, res)

	res = isOperatorChar(input[2])
	assert.Equal(t, true, res)

	res = isOperatorChar(input[3])
	assert.Equal(t, false, res)
}
