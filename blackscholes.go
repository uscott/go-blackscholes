package blackscholes

import (
	"math"

	"github.com/pkg/errors"
)

type OptionType rune

const (
	Call     = OptionType('c')
	Put      = OptionType('p')
	Straddle = OptionType('s')
)

const InvSqrt2PI float64 = 1.0 / math.Sqrt2 / math.SqrtPi

var (
	ErrNegPremium        = errors.New("Negative option premium")
	ErrNegPrice          = errors.New("Negative underlying price")
	ErrNegStrike         = errors.New("Negative strike")
	ErrNegTimeToExp      = errors.New("Negative time to expiry")
	ErrUnknownOptionType = errors.New("Unknown option type")
	ErrNilPtrArg         = errors.New("Nil pointer argument")
	ErrNoncovergence     = errors.New("Did not converge")
)

var (
	abs  func(float64) float64          = math.Abs
	exp  func(float64) float64          = math.Exp
	inf  func(int) float64              = math.Inf
	log  func(float64) float64          = math.Log
	max  func(float64, float64) float64 = math.Max
	nan  func() float64                 = math.NaN
	sqrt func(float64) float64          = math.Sqrt
)

type PriceParams struct {
	Vol          float64
	TimeToExpiry float64
	Underlying   float64
	Strike       float64
	Rate         float64
	Dividend     float64
	Type         OptionType
}

func Price(pars *PriceParams) (price float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = CheckPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	price = BSPrice(v, t, x, k, r, q, pars.Type)
	return
}

func Delta(pars *PriceParams) (delta float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = CheckPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	delta = BSDelta(v, t, x, k, r, q, pars.Type)

	return
}

func Gamma(pars *PriceParams) (gamma float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = CheckPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	gamma = BSGamma(v, t, x, k, r, q, pars.Type)
	return
}

func Vega(pars *PriceParams) (vega float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = CheckPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	vega = BSVega(v, t, x, k, r, q, pars.Type)

	return
}

func Theta(pars *PriceParams) (theta float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = CheckPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	theta = BSTheta(v, t, x, k, r, q, pars.Type)

	return
}

// AtmApprox approximates the option price when exp(-q*t)*x == exp(-r*t)*k
// v = volatility in same units as t
// t = time to expiry
// x = value of spot/underlying
// r = continuously compounded interest rate in same units as t
// q = continuous dividend yield in same units as t
// o = option type (Call, Put, Straddle)
func AtmApprox(v, t, x, q float64, o OptionType) float64 {

	switch o {
	case Call, Put:
		return exp(-q*t) * v * x * sqrt(t) * InvSqrt2PI
	case Straddle:
		return 2 * exp(-q*t) * v * x * sqrt(t) * InvSqrt2PI
	}

	return nan()
}

// CheckPriceParams checks whether t, x, k are non-negative and
// o is one of the defined Option Types
func CheckPriceParams(t, x, k float64, o OptionType) error {

	if !ValidOptionType(o) {
		return ErrUnknownOptionType
	}

	switch {
	case t < 0:
		return ErrNegTimeToExp
	case x < 0:
		return ErrNegPrice
	case k < 0:
		return ErrNegStrike
	}
	return nil
}

func GetFloatPriceParams(pars *PriceParams) (v, t, x, k, r, q float64) {
	if pars == nil {
		panic(ErrNilPtrArg)
	}
	v, t, x, k, r, q = pars.Vol, pars.TimeToExpiry,
		pars.Underlying, pars.Strike, pars.Rate, pars.Dividend
	return
}

// BSPrice returns the Black Scholes option price.
// v = volatility in same units as t
// t = time to expiry
// x = value of spot/underlying
// k = strike price
// r = continuously compounded interest rate in same units as t
// q = continuous dividend yield in same units as t
// o = option type (Call, Put, Straddle)
func BSPrice(v, t, x, k, r, q float64, o OptionType) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	return BSPriceNoErrorCheck(v, t, x, k, r, q, o)
}

func BSPriceNoErrorCheck(v, t, x, k, r, q float64, o OptionType) float64 {

	if v < 0 {
		p := BSPrice(-v, t, x, k, r, q, o)
		i := Intrinsic(t, x, k, r, q, o)
		e := p - i
		return i - e
	}

	switch {
	case x == 0:
		return ZeroUnderlyingBSPrice(t, k, r, o)
	case k == 0:
		return ZeroStrikeBSPrice(t, x, q, o)
	case v == 0:
		return Intrinsic(t, x, k, r, q, o)
	}

	d1 := D1(v, t, x, k, r, q)
	d2 := D2fromD1(d1, v, t)
	Nd1, Nd2 := NormCDF(d1), NormCDF(d2)
	dfq, dfr := exp(-q*t), exp(-r*t)

	switch o {
	case Call:
		return dfq*Nd1*x - dfr*Nd2*k
	case Put:
		return dfq*(Nd1-1)*x - dfr*(Nd2-1)*k
	}

	return dfq*(2*Nd1-1)*x - dfr*(2*Nd2-1)*k
}

func BSDelta(v, t, x, k, r, q float64, o OptionType) float64 {

	if v < 0 {
		return 2*ZeroVolBSDelta(t, x, k, r, q, o) - BSDelta(-v, t, x, k, r, q, o)
	}

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	switch {
	case t < 0, x < 0, k < 0:
		return nan()
	case x == 0:
		return ZeroUnderlyingBSDelta(t, q, o)
	case k == 0:
		return ZeroStrikeBSDelta(t, q, o)
	case v == 0:
		return ZeroVolBSDelta(t, x, k, r, q, o)
	}

	Nd1 := NormCDF(D1(v, t, x, k, r, q))

	switch o {
	case Call:
		return exp(-q*t) * Nd1
	case Put:
		return exp(-q*t) * (Nd1 - 1)
	}

	return exp(-q*t) * (2*Nd1 - 1)
}

func BSGamma(v, t, x, k, r, q float64, o OptionType) float64 {

	if v < 0 {
		return 2*ZeroVolBSGamma(t, x, k, r, q) - BSGamma(-v, t, x, k, r, q, o)
	}
	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	switch {
	case t < 0, x < 0, k < 0:
		return nan()
	case x == 0, k == 0:
		return 0
	case v == 0:
		return ZeroVolBSGamma(t, x, k, r, q)
	}

	d1 := D1(v, t, x, k, r, q)

	if o == Call || o == Put {
		return exp(-q*t-d1*d1/2) / x / v / sqrt(t) * InvSqrt2PI
	}

	return 2 * exp(-q*t-d1*d1/2) / x / v / sqrt(t) * InvSqrt2PI
}

func BSTheta(v, t, x, k, r, q float64, o OptionType) float64 {

	if v < 0 {
		return 2*ZeroVolBSTheta(t, x, k, r, q, o) - BSTheta(-v, t, x, k, r, q, o)
	}

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	switch {
	case t < 0, x < 0, k < 0:
		return nan()
	case x == 0:
		return ZeroUnderlyingBSTheta(t, k, r, o)
	case k == 0:
		return ZeroStrikeBSTheta(t, x, q, o)
	case v == 0:
		return ZeroVolBSTheta(t, x, k, r, q, o)
	case t == 0:
		return inf(-1)
	}

	d1 := D1(v, t, x, k, r, q)
	d2 := D2fromD1(d1, v, t)

	theta := -exp(-q*t) * v * x * exp(-d1*d1/2) / 2 / sqrt(t) * InvSqrt2PI
	theta += q*x*exp(-q*t)*NormCDF(d1) - r*k*exp(-r*t)*NormCDF(d2)

	if o == Call || o == Put {
		return theta
	}

	return 2 * theta
}

func BSVega(v, t, x, k, r, q float64, o OptionType) float64 {

	if v < 0 {
		return -BSVega(-v, t, x, k, r, q, o)
	}

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	if v == 0 || t == 0 || x == 0 || k == 0 {
		return 0
	}

	d1 := D1(v, t, x, k, r, q)

	if o == Call || o == Put {
		return x * exp(-q*t-d1*d1/2) * sqrt(t) * InvSqrt2PI
	}

	return 2 * x * exp(-q*t-d1*d1/2) * sqrt(t) * InvSqrt2PI
}

func D2fromD1(d1, v, t float64) float64 {
	return d1 - v*sqrt(t)
}

func D1(v, t, x, k, r, q float64) float64 {
	return (log(x/k) + (r-q+0.5*v*v)*t) / v / sqrt(t)
}

func D2(v, t, x, k, r, q float64) float64 {
	return (log(x/k) + (r-q-0.5*v*v)*t) / v / sqrt(t)
}

func Intrinsic(t, x, k, r, q float64, o OptionType) float64 {

	f := exp(-q*t)*x - exp(-r*t)*k

	switch o {
	case Call:
		return max(0, +f)
	case Put:
		return max(0, -f)
	}
	return abs(f)
}

func ValidOptionType(o OptionType) bool {
	return o == Call || o == Put || o == Straddle
}

func ZeroStrikeBSPrice(t, x, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return exp(-q*t) * x
	case Put:
		return 0
	}
	return nan()
}

func ZeroUnderlyingBSPrice(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return exp(-r*t) * k
	}
	return nan()
}

func ZeroStrikeBSDelta(t, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return exp(-q * t)
	case Put:
		return 0
	}
	return nan()
}

func ZeroUnderlyingBSDelta(t, q float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return -exp(-q * t)
	}
	return nan()
}

func ZeroVolBSDelta(t, x, k, r, q float64, o OptionType) float64 {
	q = exp(-q * t)
	x, k = q*x, exp(-r*t)*k
	switch o {
	case Call:
		if x < k {
			return 0
		}
		return q
	case Put:
		if x < k {
			return -q
		}
		return 0
	case Straddle:
		if x < k {
			return -q
		}
		return q
	}
	return nan()
}

func ZeroVolBSGamma(t, x, k, r, q float64) float64 {
	if exp(-q*t)*x != exp(-r*t)*k {
		return 0
	}
	return inf(1)
}

func ZeroStrikeBSTheta(t, x, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return q * x * exp(-q*t)
	case Put:
		return 0
	}
	return nan()
}

func ZeroUnderlyingBSTheta(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return r * k * exp(-r*t)
	}
	return nan()
}

func ZeroVolBSTheta(t, x, k, r, q float64, o OptionType) float64 {

	x, k = exp(-q*t)*x, exp(-r*t)*k

	switch o {
	case Call:
		switch {
		case x > k:
			return q*x - r*k
		case x < k:
			return 0
		default:
			return q * x
		}
	case Put:
		switch {
		case x > k:
			return 0
		case x < k:
			return r*k - q*x
		default:
			return r * k
		}
	case Straddle:
		switch {
		case x > k:
			return q*x - r*k
		case x < k:
			return r*k - q*x
		default:
			return q*x + r*k
		}
	}

	return nan()
}
