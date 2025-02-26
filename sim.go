package blackscholes

import (
	"errors"
	"math"
)

const defaultNumPaths uint = 10000000

func PriceSim(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
	numPaths ...uint,
) (price float64, err error) {

	npaths := defaultNumPaths
	if len(numPaths) > 0 {
		npaths = numPaths[0]
	}

	price = math.NaN()

	if npaths == 0 {
		err = errors.New("number of paths must be positive")
		return
	}

	if !ValidOptionType(optionType) {
		err = ErrUnknownOptionType
		return
	}

	expectedSpot := spot * math.Exp((interestRate-dividendYield)*timeToExpiry)
	sigma := vol * math.Sqrt(timeToExpiry)
	mu := -0.5 * sigma * sigma

	sum := 0.0

	for i := uint(1); i < npaths; i += 2 {

		u := float64(i) / float64(npaths)
		z := NormCDFInverse(u)

		spot = expectedSpot * math.Exp(mu+sigma*z)
		sum += Intrinsic(0, spot, strike, 0, 0, optionType)

		spot = expectedSpot * math.Exp(mu-sigma*z)
		sum += Intrinsic(0, spot, strike, 0, 0, optionType)

	}

	if npaths%2 == 1 {
		spot = expectedSpot * math.Exp(mu)
		sum += Intrinsic(0, spot, strike, 0, 0, optionType)
	}

	price = math.Exp(-interestRate*timeToExpiry) * sum / float64(npaths)

	return
}
