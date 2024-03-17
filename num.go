package blackscholes

import (
	"math"

	"github.com/pkg/errors"
)

const defaultEpsilon float64 = 1.0 / (1 << 30)

var ErrNegativeEpsilon = errors.New("epsilon must be positive")

func getEpsilon(epsilon ...float64) (eps float64, err error) {
	eps = defaultEpsilon
	if len(epsilon) > 0 {
		eps = epsilon[0]
	}
	if eps <= 0 {
		err = ErrNegativeEpsilon
	}
	return
}

func DeltaNumeric(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
	epsilon ...float64,
) (delta float64, err error) {

	delta = math.NaN()

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		return
	}

	eps, err := getEpsilon(epsilon...)
	if err != nil {
		return
	}

	upPrice, err := Price(
		vol,
		timeToExpiry,
		spot+eps,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)

	if err != nil {
		return
	}

	if spot > eps {
		spot -= eps
		eps *= 2
	}

	downPrice, err := Price(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	delta = (upPrice - downPrice) / eps
	return
}

func GammaNumeric(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
	epsilon ...float64,
) (gamma float64, err error) {

	gamma = math.NaN()

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		return
	}

	eps := defaultEpsilon
	if len(epsilon) > 0 {
		eps = epsilon[0]
	}

	if eps <= 0 {
		err = ErrNegativeEpsilon
		return
	}

	deltaUp, err := Delta(
		vol,
		timeToExpiry,
		spot+eps,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	if spot > eps {
		spot -= eps
		eps *= 2
	}

	deltaDown, err := Delta(
		vol,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	gamma = (deltaUp - deltaDown) / eps

	return
}

func VegaNumeric(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
	epsilon ...float64,
) (vega float64, err error) {

	vega = math.NaN()

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		return
	}

	eps := defaultEpsilon
	if len(epsilon) > 0 {
		eps = epsilon[0]
	}

	priceUp, err := Price(
		vol+eps,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	priceDown, err := Price(
		vol-eps,
		timeToExpiry,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	vega = 0.5 * (priceUp - priceDown) / eps

	return
}

func ThetaNumeric(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
	epsilon ...float64,
) (theta float64, err error) {

	theta = math.NaN()

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		return
	}

	eps, err := getEpsilon(epsilon...)
	if err != nil {
		return
	}

	// Note the negative sign on eps
	priceDown, err := Price(
		vol,
		timeToExpiry+eps,
		spot,
		strike,
		interestRate,
		dividendYield,
		optionType,
	)
	if err != nil {
		return
	}

	if timeToExpiry > eps {
		timeToExpiry -= eps
		eps *= 2
	}

	priceUp, err := Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	if err != nil {
		return
	}

	theta = (priceUp - priceDown) / eps

	return
}
