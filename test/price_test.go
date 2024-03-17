package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	bs "github.com/uscott/go-blackscholes"
)

const (
	numericEpsilon = 1e-32
	testEpsilon    = 1e-8
)

func getTestParams() (vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType bs.OptionType) {
	vol = 1
	timeToExpiry = 1
	spot = 100
	strike = 100
	interestRate = 0
	dividendYield = 0
	optionType = bs.Call
	return
}

func TestPrice(t *testing.T) {

	assert := assert.New(t)

	actual, err := bs.Price(0, 0, 0, 0, 0, 0, bs.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(actual))

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType := getTestParams()

	actual, err = bs.Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	expected := 7.9655792417
	assert.NoError(err)
	assert.InDelta(expected, actual, testEpsilon)

	// price1 := actual
	// price2, err := bs.PriceSim(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	// assert.NoError(err)
	// assert.InDelta(price1, price2, testEpsilon)
}
