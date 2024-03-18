package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-blackscholes"
)

func TestGamma(t *testing.T) {

	assert := assert.New(t)

	gamma, err := blackscholes.Gamma(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(gamma))

	tolerance := 1e-4
	vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType := getTestParams()

	gamma, err = blackscholes.Gamma(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	assert.NoError(err)

	gammaNum, err := blackscholes.GammaNumeric(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	assert.NoError(err)

	assert.InDelta(gamma, gammaNum, tolerance)
}
