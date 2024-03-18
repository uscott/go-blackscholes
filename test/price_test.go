package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uscott/go-blackscholes"
)

const defaultTolerance float64 = 1e-4

func getTestParams() (vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType blackscholes.OptionType) {
	vol = 0.2
	timeToExpiry = 1
	spot = 100
	strike = 100
	interestRate = 0.05
	dividendYield = 0.01
	optionType = blackscholes.Call
	return
}

func TestPrice(t *testing.T) {

	assert := assert.New(t)
	tolerance := defaultTolerance

	actual, err := blackscholes.Price(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(actual))

	vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType := getTestParams()

	expected := 9.8262858235
	actual, err = blackscholes.Price(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	assert.NoError(err)
	assert.InEpsilon(expected, actual, tolerance)

	expected = 5.9442448987
	optionType = blackscholes.Put
	actual, err = blackscholes.Price(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	assert.NoError(err)
	assert.InEpsilon(expected, actual, tolerance)

	expected = 15.7705307222
	optionType = blackscholes.Straddle
	actual, err = blackscholes.Price(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	assert.NoError(err)
	assert.InEpsilon(expected, actual, tolerance)

	tolerance = 1e-3
	price1 := actual
	price2, err := blackscholes.PriceSim(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	assert.NoError(err)
	assert.InEpsilon(price1, price2, tolerance)
}
