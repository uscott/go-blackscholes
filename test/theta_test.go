package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-blackscholes"
)

func TestTheta(t *testing.T) {

	assert := assert.New(t)

	theta, err := blackscholes.Theta(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(theta))
}
