package blackscholes

import (
	"math"

	"github.com/pkg/errors"
)

const defaultEpsilon float64 = 1.0 / (1 << 30)

var ErrNegativeEpsilon = errors.New("epsilon must be positive")

func DeltaNumeric(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType, epsilon ...float64) (delta float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		delta = math.NaN()
		return
	}

	eps := defaultEpsilon

	if len(epsilon) > 0 {
		eps = epsilon[0]
	}

	if eps <= 0 {
		delta = math.NaN()
		err = ErrNegativeEpsilon
		return
	}

	var upPrice, downPrice float64

	upPrice, err = Price(vol, timeToExpiry, spot+eps, strike, interestRate, dividendYield, optionType)

	if spot < eps {
		var midPrice float64
		midPrice, err = Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		if err != nil {
			delta = math.NaN()
			return
		}
		downPrice = ZeroUnderlyingBSPrice(timeToExpiry, strike, interestRate, optionType)
		upWeight := spot / eps / (spot + eps)
		midWeight := (eps - spot) / eps / (spot + eps)
		downWeight := -eps / spot / (spot + eps)
		delta = upWeight*upPrice + midWeight*midPrice + downWeight*downPrice
		return
	}

	downPrice, err = Price(vol, timeToExpiry, spot-eps, strike, interestRate, dividendYield, optionType)
	if err != nil {
		delta = math.NaN()
		return
	}

	delta = (upPrice - downPrice) / 2 / eps

	return
}

func GammaNumeric(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType, epsilon ...float64) (gamma float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		gamma = math.NaN()
		return
	}

	eps := defaultEpsilon

	if len(epsilon) > 0 {
		eps = epsilon[0]
	}

	if eps <= 0 {
		gamma = math.NaN()
		err = ErrNegativeEpsilon
		return
	}
	if CheckPriceParams(t, x, k, o) != nil {
		return math.NaN()
	}

	e := abs(eps)

	if x < e {
		pu := BSPriceNoErrorCheck(v, t, x+e, k, r, q, o)
		px := BSPriceNoErrorCheck(v, t, 2*x, k, r, q, o)
		pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)
		pd := ZeroUnderlyingBSPrice(t, k, r, o)
		cu, cx := 2*x/(e*e)/(e+x), 2*(e-x)/e/(x*x)
		cm := -2 * (x*x*x + 2*e*e*e) / (x * x) / (e * e) / (e + x)
		cd := 2 * (e*e + x*x) / (x * x) / (e * e) / (e + x)
		return cu*pu + cx*px + cm*pm + cd*pd
	}
	pu := BSPriceNoErrorCheck(v, t, x+e, k, r, q, o)
	pd := BSPriceNoErrorCheck(v, t, x-e, k, r, q, o)
	pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)

	return (pu - 2*pm + pd) / (e * e)
}

func BSVegaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return math.NaN()
	}

	pu := BSPriceNoErrorCheck(v+eps, t, x, k, r, q, o)
	pd := BSPriceNoErrorCheck(v-eps, t, x, k, r, q, o)

	return (pu - pd) / 2 / eps
}

func BSThetaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return math.NaN()
	}

	e := abs(eps)

	if t < e {
		pu := Intrinsic(0, x, k, 0, 0, o)
		pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)
		pd := BSPriceNoErrorCheck(v, t+e, x, k, r, q, o)
		cu, cm, cd := x/e/(e+x), (e-x)/e/x, -e/x/(e+x)
		return cu*pu + cm*pm + cd*pd
	}

	pu := BSPriceNoErrorCheck(v, t-e, x, k, r, q, o)
	pd := BSPriceNoErrorCheck(v, t+e, x, k, r, q, o)

	return (pu - pd) / 2 / e
}
