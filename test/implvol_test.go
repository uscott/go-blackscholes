package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-blackscholes"
)

func TestImpliedVol(t *testing.T) {

	assert := assert.New(t)

	vol, err := blackscholes.ImpliedVol(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(vol))

	vol, err = blackscholes.ImpliedVol(0, 0, 0, 0, 0, 0, blackscholes.Straddle)
	tolerance := 1e-4
	assert.NoError(err)
	assert.InDelta(0, vol, tolerance)
}
