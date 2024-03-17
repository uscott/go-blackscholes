package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bs "github.com/uscott/go-blackscholes"
)

func Test_Gamma(t *testing.T) {

	assert := assert.New(t)

	gamma, err := bs.Gamma(0, 0, 0, 0, 0, 0, bs.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(gamma))

	tolerance := 1e-4
	vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType := getTestParams()

	gamma, err = bs.Gamma(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	assert.NoError(err)

	gammaNum, err := bs.GammaNumeric(
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
