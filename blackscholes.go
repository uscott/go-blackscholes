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

func ValidOptionType(optionType OptionType) bool {
	return optionType == Call || optionType == Put || optionType == Straddle
}

func getd1d2(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
) (d1, d2 float64, err error) {

	d1 = math.NaN()
	d2 = math.NaN()

	if timeToExpiry <= 0 {
		err = ErrNegTimeToExp
		return
	}

	if vol <= 0 {
		err = ErrNegVol
		return
	}

	if spot <= 0 || strike <= 0 {
		err = fmt.Errorf("invalid spot (%f) or strike price (%f)", spot, strike)
		return
	}

	spot *= math.Exp(-dividendYield * timeToExpiry)
	strike *= math.Exp(-interestRate * timeToExpiry)
	vol *= math.Sqrt(timeToExpiry)

	d1 = math.Log(spot/strike)/vol + 0.5*vol
	d2 = d1 - vol

	return
}

// Price returns the Black Scholes option price.
// vol = volatility in same units as timeToExpiry
// timeToExpiry = time to expiry
// spot = value of spot/underlying
// strike = strike price
// interestRate = continuously compounded interest rate in same units as t
// dividendYield = continuous dividend yield in same units as t
// optionType = option type (Call, Put, Straddle)

func Price(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) (price float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		price = math.NaN()
		return
	}

	switch {
	case spot == 0:
		price = PriceZeroSpot(timeToExpiry, strike, interestRate, optionType)
		return
	case strike == 0:
		price = PriceZeroStrike(timeToExpiry, spot, dividendYield, optionType)
		return
	case vol == 0:
		price = Intrinsic(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	}

	volIsNegative := vol < 0

	vol = math.Abs(vol)

	d1, d2, err := getd1d2(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
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

func Delta(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) (delta float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		delta = math.NaN()
		return
	}

	switch {
	case spot == 0:
		delta = DeltaZeroSpot(timeToExpiry, dividendYield, optionType)
		return
	case strike == 0:
		delta = DeltaZeroStrike(timeToExpiry, spot, optionType)
		return
	case vol == 0:
		delta = DeltaZeroVol(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	d1, _, err := getd1d2(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
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
		zeroVolDelta := DeltaZeroVol(
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		)
		delta = zeroVolDelta - (delta - zeroVolDelta)
	}

	return
}

// AtmApprox approximates the option price when the spot discounted by the dividend yield
// is equal to the strike discounted by the interest rate.
func AtmApprox(
	vol, timeToExpiry, spot, dividendYield float64,
	optionType OptionType,
) (price float64, err error) {

	price = math.NaN()

	if timeToExpiry < 0 {
		err = ErrNegTimeToExp
		return
	}

	if !ValidOptionType(optionType) {
		err = ErrUnknownOptionType
		return
	}

	price = math.Exp(
		-dividendYield*timeToExpiry,
	) * spot * vol * math.Sqrt(
		timeToExpiry,
	) * InvSqrt2PI

	if optionType == Straddle {
		price *= 2
	}

	return
}

func Gamma(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) (gamma float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		gamma = math.NaN()
		return
	}

	switch {
	case spot == 0, strike == 0:
		gamma = 0
		return
	case vol == 0:
		gamma = GammaZeroVol(timeToExpiry, spot, strike, interestRate, dividendYield)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	d1, _, err := getd1d2(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		gamma = math.NaN()
		return
	}

	gamma = math.Exp(
		-dividendYield*timeToExpiry,
	) * NormPDF(
		d1,
	) / (spot * vol * math.Sqrt(timeToExpiry))

	if optionType == Straddle {
		gamma *= 2
	}

	if volIsNegative {
		gamma = 2*GammaZeroVol(timeToExpiry, spot, strike, interestRate, dividendYield) - gamma
	}

	return
}

func Theta(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) (theta float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		theta = math.NaN()
		return
	}

	switch {
	case spot == 0:
		theta = ThetaZeroSpot(timeToExpiry, strike, interestRate, optionType)
		return
	case strike == 0:
		theta = ThetaZeroStrike(timeToExpiry, spot, dividendYield, optionType)
		return
	case vol == 0:
		theta = ThetaZeroVol(timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
		return
	case timeToExpiry == 0:
		theta = math.Inf(-1)
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	d1, d2, err := getd1d2(vol, timeToExpiry, spot, strike, interestRate, dividendYield)
	if err != nil {
		theta = math.NaN()
		return
	}

	spot *= math.Exp(-dividendYield * timeToExpiry)
	strike *= math.Exp(-interestRate * timeToExpiry)

	switch optionType {
	case Call:
		theta = -0.5*vol*spot*NormPDF(
			d1,
		)/math.Sqrt(
			timeToExpiry,
		) - interestRate*strike*NormCDF(
			d2,
		) + dividendYield*spot*NormCDF(
			d1,
		)
	case Put:
		theta = -0.5*vol*spot*NormPDF(
			d1,
		)*math.Sqrt(
			timeToExpiry,
		) + interestRate*strike*NormCDF(
			-d2,
		) - dividendYield*spot*NormCDF(
			-d1,
		)
	case Straddle:
		theta = -vol*spot*NormPDF(
			d1,
		)/math.Sqrt(
			timeToExpiry,
		) - interestRate*strike*(NormCDF(d2)-NormCDF(-d2)) + dividendYield*spot*(NormCDF(d1)-NormCDF(-d1))
	}

	if volIsNegative {
		theta = 2*ThetaZeroVol(
			timeToExpiry,
			spot,
			strike,
			interestRate,
			dividendYield,
			optionType,
		) - theta
	}

	return
}

func Vega(
	vol, timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) (vega float64, err error) {

	if err = CheckPriceParams(timeToExpiry, spot, strike, optionType); err != nil {
		vega = math.NaN()
		return
	}

	if vol == 0 || timeToExpiry == 0 || spot == 0 || strike == 0 {
		return
	}

	volIsNegative := vol < 0
	vol = math.Abs(vol)

	d1, _, err := getd1d2(vol, timeToExpiry, spot, strike, interestRate, dividendYield)

	vega = spot * math.Exp(
		-dividendYield*timeToExpiry-0.5*d1*d1,
	) * math.Sqrt(
		timeToExpiry,
	) * InvSqrt2PI

	if optionType == Straddle {
		vega *= 2
	}

	if volIsNegative {
		vega *= -1
	}

	return
}

func Intrinsic(
	timeToExpiry, spot, strike, interestRate, dividendYield float64,
	optionType OptionType,
) float64 {

	forwardValue := math.Exp(
		-dividendYield*timeToExpiry,
	)*spot - math.Exp(
		-interestRate*timeToExpiry,
	)*strike

	switch optionType {
	case Call:
		return math.Max(0, forwardValue)
	case Put:
		return math.Max(0, -forwardValue)
	}
	return math.Abs(forwardValue)
}

func PriceZeroStrike(t, x, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return math.Exp(-q*t) * x
	case Put:
		return 0
	}
	return math.NaN()
}

func PriceZeroSpot(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return math.Exp(-r*t) * k
	}
	return math.NaN()
}

func DeltaZeroStrike(t, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return math.Exp(-q * t)
	case Put:
		return 0
	}
	return math.NaN()
}

func DeltaZeroSpot(t, q float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return -math.Exp(-q * t)
	}
	return math.NaN()
}

func DeltaZeroVol(t, x, k, r, q float64, o OptionType) float64 {

	dfq := math.Exp(-q * t)
	x, k = dfq*x, math.Exp(-r*t)*k

	switch o {
	case Call:
		if x < k {
			return 0
		}
		if x > k {
			return dfq
		}
		return 0.5 * dfq // Convention to match numeric delta
	case Put:
		if x < k {
			return -dfq
		}
		if x > k {
			return 0
		}
		return -0.5 * dfq // Convention to match numeric delta
	case Straddle:
		if x < k {
			return -dfq
		}
		if x > k {
			return dfq
		}
		return 0 // Convention to match numeric delta
	}
	return math.NaN()
}

func GammaZeroVol(t, x, k, r, q float64) float64 {
	if math.Exp(-q*t)*x-math.Exp(-r*t)*k != 0 {
		return 0
	}
	return math.Inf(1)
}

func ThetaZeroStrike(t, x, q float64, o OptionType) float64 {
	switch o {
	case Call, Straddle:
		return q * x * math.Exp(-q*t)
	case Put:
		return 0
	}
	return math.NaN()
}

func ThetaZeroSpot(t, k, r float64, o OptionType) float64 {
	switch o {
	case Call:
		return 0
	case Put, Straddle:
		return r * k * math.Exp(-r*t)
	}
	return math.NaN()
}

func ThetaZeroVol(t, x, k, r, q float64, o OptionType) float64 {

	x, k = math.Exp(-q*t)*x, math.Exp(-r*t)*k

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
