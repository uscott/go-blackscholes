package blackscholes

import (
	"errors"
	"fmt"
	"math"
)

const (
	defaultTolerance     float64 = 1.0 / (1 << 30)
	defaultUpperBound    float64 = 0.01
	defaultLowerBound    float64 = 1.99
	defaultMaxIterations int     = 1000000
)

var ErrMaxIterations = errors.New("max iterations exceeded")

type ImpliedVolParams struct {
	LowerBound    *float64
	UpperBound    *float64
	Tolerance     *float64
	MaxIterations *int
}

func ImpliedVol(premium, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType, params ...ImpliedVolParams) (vol float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		vol = math.NaN()
		return
	}

	intrinsic := Intrinsic(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	extrinsic := premium - intrinsic

	if math.Abs(extrinsic) <= math.SmallestNonzeroFloat64 {
		return
	}

	tol, lb, ub, maxit := defaultTolerance, defaultLowerBound, defaultUpperBound, defaultMaxIterations

	if len(params) > 0 {
		p := params[0]
		if p.LowerBound != nil {
			lb = *p.LowerBound
		}
		if p.UpperBound != nil {
			ub = *p.UpperBound
		}
		if p.Tolerance != nil {
			tol = *p.Tolerance
		}
		if p.MaxIterations != nil {
			maxit = *p.MaxIterations
		}
	}

	vol = math.NaN()
	var lowPrice, highPrice float64

	lowPrice, err = Price(lb, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	if err != nil {
		return
	}

	highPrice, err = Price(ub, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
	if err != nil {
		return
	}

	var it int
	const boundIncrement float64 = 0.47

	// Adjust bounds
	for ; premium < lowPrice || highPrice < premium; it++ {
		if it > maxit {
			err = ErrMaxIterations
			return
		}
		if premium < lowPrice {
			lb -= boundIncrement
			lowPrice, err = Price(lb, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
			if err != nil {
				return
			}
		}
		if highPrice < premium {
			ub += boundIncrement
			highPrice, err = Price(ub, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
			if err != nil {
				return
			}
		}
	}

	if premium < lowPrice {
		err = fmt.Errorf("failed to find lower bound - lower bound price, lower bound vol, iterations: %v, %v, %d", lowPrice, lb, it)
		return
	}

	if highPrice < premium {
		err = fmt.Errorf("failed to find upper bound - upper bound price, upper bound vol, iterations: %v, %v, %d", lowPrice, lb, it)
		return
	}

	// Bisection
	var price float64

	for ; ub-lb > tol; it++ {
		if it > maxit {
			vol = math.NaN()
			err = ErrMaxIterations
			return
		}
		vol = 0.5 * (lb + ub)
		price, err = Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		if err != nil {
			vol = math.NaN()
			return
		}
		if premium < price {
			ub = vol
		} else {
			lb = vol
		}
	}

	vol = CorrectVolSign(extrinsic, vol)

	return
}

func CorrectVolSign(extrinsic float64, vol float64) float64 {
	if extrinsic > 0 && vol < 0 || extrinsic < 0 && vol > 0 {
		return -vol
	}
	return vol
}
