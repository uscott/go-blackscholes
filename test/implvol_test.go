package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bs "github.com/uscott/go-blackscholes"
)

func Test_ImpliedVol(t *testing.T) {

	assert := assert.New(t)

	vol, err := bs.ImpliedVol(0, 0, 0, 0, 0, 0, bs.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(vol))

	vol, err = bs.ImpliedVol(0, 0, 0, 0, 0, 0, bs.Straddle)
	assert.NoError(err)
	assert.InDelta(0, vol, testEpsilon)
}
