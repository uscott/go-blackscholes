package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	blackscholes "github.com/uscott/go-blackscholes"
)

func TestImpliedVol(t *testing.T) {

	assert := assert.New(t)
	tolerance := defaultTolerance

	vol, err := blackscholes.ImpliedVol(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(vol))

	vol, err = blackscholes.ImpliedVol(0, 0, 0, 0, 0, 0, blackscholes.Straddle)
	assert.NoError(err)
	assert.InDelta(0, vol, tolerance)

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []blackscholes.OptionType{blackscholes.Call, blackscholes.Put, blackscholes.Straddle} {
		premium, err := blackscholes.Price(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		impliedVol, err := blackscholes.ImpliedVol(
			premium,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		assert.InDelta(vol, impliedVol, tolerance)
	}
}
