package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	blackscholes "github.com/uscott/go-blackscholes"
)

func TestTheta(t *testing.T) {

	t.Skip("Need to fix theta")

	assert := assert.New(t)
	tolerance := defaultTolerance

	theta, err := blackscholes.Theta(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(theta))

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, _ := getTestParams()

	for _, optionType := range []blackscholes.OptionType{blackscholes.Call, blackscholes.Put, blackscholes.Straddle} {
		theta, err = blackscholes.Theta(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)

		thetaNum, err := blackscholes.ThetaNumeric(
			vol,
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		assert.NoError(err)
		assert.InDelta(theta, thetaNum, tolerance)
	}
}
