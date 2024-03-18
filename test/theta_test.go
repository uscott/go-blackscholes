package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bs "github.com/uscott/go-blackscholes"
)

func TestTheta(t *testing.T) {

	assert := assert.New(t)

	theta, err := bs.Theta(0, 0, 0, 0, 0, 0, bs.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(theta))
}
