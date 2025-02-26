package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	blackscholes "github.com/uscott/go-blackscholes"
)

func TestDelta(t *testing.T) {

	assert := assert.New(t)
	tolerance := defaultTolerance

	delta, err := blackscholes.Delta(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(delta))

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []blackscholes.OptionType{blackscholes.Call, blackscholes.Put, blackscholes.Straddle} {
		delta, err = blackscholes.Delta(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		deltaNum, err := blackscholes.DeltaNumeric(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		assert.InDelta(delta, deltaNum, tolerance)
	}
}

func TestDeltaZeroVol(t *testing.T) {

	assert := assert.New(t)
	tolerance := defaultTolerance

	_, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []blackscholes.OptionType{blackscholes.Call, blackscholes.Put, blackscholes.Straddle} {
		delta := blackscholes.DeltaZeroVol(
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		deltaNum, err := blackscholes.DeltaNumeric(
			0,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		assert.InDelta(delta, deltaNum, tolerance)
	}
}
