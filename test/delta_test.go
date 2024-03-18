package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bs "github.com/uscott/go-blackscholes"
)

func TestDelta(t *testing.T) {

	assert := assert.New(t)

	delta, err := bs.Delta(0, 0, 0, 0, 0, 0, bs.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(delta))

	tolerance := 1e-4
	vol, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []bs.OptionType{bs.Call, bs.Put, bs.Straddle} {
		delta, err = bs.Delta(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		deltaNum, err := bs.DeltaNumeric(
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
