package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-blackscholes"
)

func TestVega(t *testing.T) {

	assert := assert.New(t)
	tolerance := defaultTolerance

	vega, err := blackscholes.Vega(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(vega))

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []blackscholes.OptionType{blackscholes.Call, blackscholes.Put, blackscholes.Straddle} {
		vega, err = blackscholes.Vega(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)

		vegaNum, err := blackscholes.VegaNumeric(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		assert.InDelta(vega, vegaNum, tolerance)
	}
}
