package blackscholes

import (
	"fmt"
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
	ErrNegVol            = errors.New("Negative volatility")
	ErrNegTimeToExp      = errors.New("Negative time to expiry")
	ErrNegPremium        = errors.New("Negative option premium")
	ErrNegPrice          = errors.New("Negative underlying price")
	ErrNegStrike         = errors.New("Negative strike")
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

// Price returns the Black Scholes option price.
// vol = volatility in same units as timeToExpiry
// timeToExpiry = time to expiry
// spot = value of spot/underlying
// strike = strike price
// interestRate = continuously compounded interest rate in same units as t
// dividendYield = continuous dividend yield in same units as t
// optionType = option type (Call, Put, Straddle)

func Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) (price float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		price = math.NaN()
		return
	}

	switch {
	case spot == 0:
		price = ZeroUnderlyingBSPrice(timeToExpiry, strike, interestRate, optionType)
		return
	case strike == 0:
		price = ZeroStrikeBSPrice(timeToExpiry, spot, dividendYield, optionType)
		return
	case vol == 0:
		price = Intrinsic(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	}

	volIsNegative := vol < 0

	vol = math.Abs(vol)

	var d1, d2 float64

	d1, err = D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		price = math.NaN()
		return
	}

	d2, err = D2fromD1(d1, vol, timeToExpiry)
	if err != nil {
		price = math.NaN()
		return
	}

	Nd1, Nd2 := NormCDF(d1), NormCDF(d2)

	spot *= math.Exp(-dividendYield * timeToExpiry)
	strike *= math.Exp(-interestRate * timeToExpiry)

	switch optionType {
	case Call:
		price = spot*Nd1 - strike*Nd2
	case Put:
		price = spot*(Nd1-1) - strike*(Nd2-1)
	case Straddle:
		price = (2*Nd1-1)*spot - (2*Nd2-1)*strike
	}

	if volIsNegative {
		intrinsic := Intrinsic(timeToExpiry, spot, strike, 0, 0, optionType)
		price = intrinsic - (price - intrinsic)
	}

	return
}

func Delta(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) (delta float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		delta = math.NaN()
		return
	}

	switch {
	case spot == 0:
		delta = ZeroUnderlyingBSDelta(timeToExpiry, dividendYield, optionType)
		return
	case strike == 0:
		delta = ZeroStrikeBSDelta(timeToExpiry, spot, optionType)
		return
	case vol == 0:
		delta = ZeroVolBSDelta(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	var d1 float64
	d1, err = D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		delta = math.NaN()
		return
	}

	Nd1 := NormCDF(d1)
	dividendDiscount := math.Exp(-dividendYield * timeToExpiry)

	switch optionType {
	case Call:
		delta = dividendDiscount * Nd1
	case Put:
		delta = dividendDiscount * (Nd1 - 1)
	case Straddle:
		delta = dividendDiscount * (2*Nd1 - 1)
	default:
		delta = math.NaN()
		err = ErrUnknownOptionType
		return
	}

	if volIsNegative {
		zeroVolDelta := ZeroVolBSDelta(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		delta = zeroVolDelta - (delta - zeroVolDelta)
	}

	return
}

// AtmApprox approximates the option price when the spot discounted by the dividend yield
// is equal to the strike discounted by the interest rate.
func AtmApprox(vol, timeToExpiry, spot, dividendYield float64, optionType OptionType) (price float64, err error) {

	if timeToExpiry < 0 {
		price = math.NaN()
		err = ErrNegTimeToExp
		return
	}

	vol *= math.Sqrt(timeToExpiry)
	spot *= math.Exp(-dividendYield * timeToExpiry)

	switch optionType {
	case Call, Put:
		price = spot * vol * InvSqrt2PI
	case Straddle:
		price = 2 * spot * vol * InvSqrt2PI
	default:
		price = math.NaN()
		err = ErrUnknownOptionType
	}

	return
}

// CheckPriceParams checks whether timeToExpiry, spot, and strike are non-negative, and
// optionType is one of the defined OptionType constants
func CheckPriceParams(timeToExpiry, spot, strike float64, optionType OptionType) error {

	if !ValidOptionType(optionType) {
		return ErrUnknownOptionType
	}

	switch {
	case timeToExpiry < 0:
		return ErrNegTimeToExp
	case spot < 0:
		return ErrNegPrice
	case strike < 0:
		return ErrNegStrike
	}
	return nil
}

func Gamma(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) (gamma float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		gamma = math.NaN()
		return
	}

	switch {
	case spot == 0, strike == 0:
		gamma = 0
		return
	case vol == 0:
		gamma = ZeroVolBSGamma(timeToExpiry, spot, strike, interestRate, dividendYield)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	var d1 float64
	d1, err = D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		gamma = math.NaN()
		return
	}

	gamma = exp(-dividendYield*timeToExpiry-d1*d1/2) / spot / vol / sqrt(timeToExpiry) * InvSqrt2PI

	if optionType == Straddle {
		gamma *= 2
	}

	if volIsNegative {
		gamma = 2*ZeroVolBSGamma(timeToExpiry, spot, strike, interestRate, dividendYield) - gamma
	}

	return
}

func Theta(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) (theta float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		theta = math.NaN()
		return
	}

	switch {
	case spot == 0:
		theta = ZeroUnderlyingBSTheta(timeToExpiry, strike, interestRate, optionType)
		return
	case strike == 0:
		theta = ZeroStrikeBSTheta(timeToExpiry, spot, dividendYield, optionType)
		return
	case vol == 0:
		theta = ZeroVolBSTheta(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	case timeToExpiry == 0:
		theta = math.Inf(-1)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	var d1, d2 float64

	d1, err = D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		theta = math.NaN()
		return
	}

	d2, err = D2fromD1(d1, vol, timeToExpiry)
	if err != nil {
		theta = math.NaN()
		return
	}

	spot *= exp(-dividendYield * timeToExpiry)
	strike *= exp(-interestRate * timeToExpiry)
	theta = -vol*spot*exp(-d1*d1/2)/2/sqrt(timeToExpiry)*InvSqrt2PI + dividendYield*spot*NormCDF(d1) - interestRate*strike*NormCDF(d2)

	if optionType == Straddle {
		theta *= 2
	}

	if volIsNegative {
		theta = 2*ZeroVolBSTheta(timeToExpiry, spot, strike, interestRate, dividendYield, optionType) - theta
	}

	return
}

func Vega(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) (vega float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		vega = math.NaN()
		return
	}

	if vol == 0 || timeToExpiry == 0 || spot == 0 || strike == 0 {
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	var d1 float64

	d1, err = D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield)

	vega = spot * exp(-dividendYield*timeToExpiry-0.5*d1*d1) * sqrt(timeToExpiry) * InvSqrt2PI

	if optionType == Straddle {
		vega *= 2
	}

	if volIsNegative {
		vega *= -1
	}

	return
}

func D2fromD1(d1, vol, timeToExpiry float64) (d2 float64, err error) {

	if timeToExpiry < 0 {
		err = ErrNegTimeToExp
		d2 = math.NaN()
		return
	}

	d2 = d1 - vol*sqrt(timeToExpiry)
	return
}

func D1(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64) (d1 float64, err error) {

	d1 = math.NaN()

	if timeToExpiry < 0 {
		err = ErrNegTimeToExp
		return
	}

	if vol < 0 {
		err = ErrNegVol
		return
	}

	if strike == 0 || spot*strike < 0 {
		err = fmt.Errorf("invalid spot (%f) or strike price (%f)", spot, strike)
		return
	}

	d1 = (math.Log(spot/strike) + (interestRate-dividendYield+0.5*vol*vol)*timeToExpiry) / vol / sqrt(timeToExpiry)

	return
}

func D2(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64) (d2 float64, err error) {

	d2 = math.NaN()

	if timeToExpiry < 0 {
		err = ErrNegTimeToExp
		return
	}

	if vol < 0 {
		err = ErrNegVol
		return
	}

	if strike == 0 || spot*strike < 0 {
		err = fmt.Errorf("invalid spot (%f) or strike price (%f)", spot, strike)
		return
	}

	d2 = (math.Log(spot/strike) + (interestRate-dividendYield-0.5*vol*vol)*timeToExpiry) / vol / sqrt(timeToExpiry)

	return
}

func Intrinsic(timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType) float64 {

	forwardValue := math.Exp(-dividendYield*timeToExpiry)*spot - math.Exp(-interestRate*timeToExpiry)*strike

	switch optionType {
	case Call:
		return math.Max(0, forwardValue)
	case Put:
		return math.Max(0, -forwardValue)
	}
	return math.Abs(forwardValue)
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
	return math.NaN()
}

func ZeroUnderlyingBSPrice(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return exp(-r*t) * k
	}
	return math.NaN()
}

func ZeroStrikeBSDelta(t, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return exp(-q * t)
	case Put:
		return 0
	}
	return math.NaN()
}

func ZeroUnderlyingBSDelta(t, q float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return -exp(-q * t)
	}
	return math.NaN()
}

func ZeroVolBSDelta(t, x, k, r, q float64, o OptionType) float64 {

	dfq := exp(-q * t)
	x, k = dfq*x, exp(-r*t)*k

	switch o {
	case Call:
		if x < k {
			return 0
		}
		return dfq
	case Put:
		if x < k {
			return -dfq
		}
		return 0
	case Straddle:
		if x < k {
			return -dfq
		}
		return dfq
	}
	return math.NaN()
}

func ZeroVolBSGamma(t, x, k, r, q float64) float64 {
	if exp(-q*t)*x-exp(-r*t)*k != 0 {
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
	return math.NaN()
}

func ZeroUnderlyingBSTheta(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return r * k * exp(-r*t)
	}
	return math.NaN()
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

	return math.NaN()
}
