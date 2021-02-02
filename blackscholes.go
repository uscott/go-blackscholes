package blackscholes

import (
	"fmt"
	"math"
	"time"

	"github.com/maurodelazeri/gaussian-distribution"
	"github.com/pkg/errors"
)

type OptionType rune

const (
	Call     = OptionType('c')
	Put      = OptionType('p')
	Straddle = OptionType('s')
)

const (
	OneMinute    float64 = 1.0 / 60 / 24 / 365
	Tiny         float64 = 1e-4
	lbDefault    float64 = 0.01
	ubDefault    float64 = 1.99
	MaxItDefault int     = 100000
	InvSqrt2PI   float64 = 1.0 / math.Sqrt2 / math.SqrtPi
)

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
	exp  func(float64) float64          = math.Exp
	log  func(float64) float64          = math.Log
	max  func(float64, float64) float64 = math.Max
	nan  func() float64                 = math.NaN
	sqrt func(float64) float64          = math.Sqrt
)

type PriceParams struct {
	Vol        float64
	Now        time.Time
	Expiry     time.Time
	Underlying float64
	Strike     float64
	Rate       float64
	Dividend   float64
	Type       OptionType
}

type ImpliedVolParams struct {
	Premium    float64
	Now        time.Time
	Expiry     time.Time
	Underlying float64
	Strike     float64
	Rate       float64
	Dividend   float64
	Type       OptionType
	LB         *float64
	UB         *float64
	Tol        *float64
	MaxIt      *int
}

func Price(pars *PriceParams) (price float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	x, k = max(x, Tiny), max(k, Tiny)
	price = bsprice(v, t, x, k, r, q, pars.Type)
	return
}

func Delta(pars *PriceParams) (delta float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	v, t, x, k = max(v, Tiny),
		max(t, OneMinute), max(x, Tiny), max(k, Tiny)

	delta = bsdelta(v, t, x, k, r, q, pars.Type)

	return
}

func Gamma(pars *PriceParams) (gamma float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	v, t, x, k = max(v, Tiny),
		max(t, OneMinute), max(x, Tiny), max(k, Tiny)

	gamma = bsgamma(v, t, x, k, r, q, pars.Type)
	return
}

func Vega(pars *PriceParams) (vega float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	v, t, x, k = max(v, Tiny),
		max(t, OneMinute), max(x, Tiny), max(k, Tiny)

	o := pars.Type

	vega = bsvega(v, t, x, k, r, q, o)

	return
}

func Theta(pars *PriceParams) (theta float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	v, t, x, k = max(v, Tiny),
		max(t, OneMinute), max(x, Tiny), max(k, Tiny)

	o := pars.Type

	theta = bstheta(v, t, x, k, r, q, o)

	return
}

func ThetaNum(pars *PriceParams) (theta float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	v, t, x, k, r, q := GetFloatPriceParams(pars)

	if err = checkPriceParams(t, x, k, pars.Type); err != nil {
		return nan(), err
	}

	v, t, x, k = max(v, Tiny),
		max(t, OneMinute), max(x, Tiny), max(k, Tiny)

	o := pars.Type

	eps := t * 365 * Tiny
	tu, td := t-eps, t+eps

	theta = bsprice(v, tu, x, k, r, q, o) - bsprice(v, td, x, k, r, q, o)
	theta /= 2 * eps

	return
}

func ImpliedVol(pars *ImpliedVolParams) (vol float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	p, t, x, k, r, q := getFloatVolParams(pars)
	o := pars.Type
	if err = checkPriceParams(t, x, k, o); err != nil {
		return nan(), err
	}
	if t < OneMinute {
		return 0, nil
	}

	x, k = max(x, Tiny), max(k, Tiny)
	lb, ub, tol, maxit := getVolSearchParams(pars)

	checkVolSearchParams(&lb, &ub, &tol, &maxit)

	var (
		it             int
		plo, phi, pmid float64
	)
	for it = 0; it < maxit; it++ {
		plo = bsprice(lb, t, x, k, r, q, o)
		if plo <= p {
			break
		}
		lb -= 0.47
	}
	if p < plo {
		return nan(), fmt.Errorf(
			"Failed to find lower bound - lb price, lb vol, iters: %v, %v, %d",
			plo, lb, it,
		)
	}
	for it = 0; it < maxit; it++ {
		phi = bsprice(ub, t, x, k, r, q, o)
		if p <= phi {
			break
		}
		ub += 0.47
	}
	if phi < p {
		return nan(), fmt.Errorf(
			"Failed to find upper bound - uprice, uvol, iters: %v, %v, %d",
			phi, ub, it,
		)
	}

	intrval := intrinsic(t, x, k, r, q, o)

	for it = 0; it < maxit; it++ {
		vol = 0.5 * (lb + ub)
		pmid = bsprice(vol, t, x, k, r, q, o)
		if ub-lb < tol {
			extrinsic := pmid - intrval
			correctVolSign(extrinsic, &vol)
			return
		}
		if p < pmid {
			ub = vol
		} else if pmid < p {
			lb = vol
		} else {
			extrinsic := pmid - intrval
			correctVolSign(extrinsic, &vol)
			return
		}
	}
	plo, phi = bsprice(lb, t, x, k, r, q, o), bsprice(ub, t, x, k, r, q, o)
	return nan(), fmt.Errorf(
		"Did not converge - lb, ub, lb price, ub price, mid, iters: %v, %v, %v, %v, %v, %d",
		lb, ub, plo, phi, pmid, it,
	)
}

func correctVolSign(extrinsic float64, vol *float64) {
	if extrinsic >= 0 && *vol < 0 {
		*vol = 0
	} else if extrinsic < 0 && *vol > 0 {
		*vol = 0
	}
}

func GetFloatPriceParams(pars *PriceParams) (v, t, x, k, r, q float64) {
	t = pars.Expiry.Sub(pars.Now).Hours() / 24 / 365
	v, x, k, r, q = pars.Vol,
		pars.Underlying, pars.Strike, pars.Rate, pars.Dividend
	return
}

func getFloatVolParams(pars *ImpliedVolParams) (p, t, x, k, r, q float64) {
	t = pars.Expiry.Sub(pars.Now).Hours() / 24 / 365
	p, x, k, r, q = pars.Premium,
		pars.Underlying, pars.Strike, pars.Rate, pars.Dividend
	return
}

func getVolSearchParams(pars *ImpliedVolParams) (lb, ub, tol float64, maxit int) {

	lb, ub, tol, maxit = lbDefault, ubDefault, Tiny, MaxItDefault
	if pars.LB != nil {
		lb = *pars.LB
	}
	if pars.UB != nil {
		ub = *pars.UB
	}
	if pars.Tol != nil {
		tol = *pars.Tol
	}
	if pars.MaxIt != nil {
		maxit = *pars.MaxIt
	}
	return
}

func checkPriceParams(t, x, k float64, o OptionType) error {

	switch o {
	case Call, Put, Straddle:
	default:
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

func checkVolSearchParams(lb, ub, tol *float64, maxit *int) {

	if lb == nil || ub == nil || tol == nil || maxit == nil {
		panic("Nil pointer")
	}

	if *ub < *lb {
		*ub = max(*lb+1, ubDefault)
	}
	if *tol <= 0 {
		*tol = Tiny
	}
	if *maxit <= 0 {
		*maxit = MaxItDefault
	}
}

func bsprice(v, t, x, k, r, q float64, o OptionType) (price float64) {

	if t < 0 || x <= 0 || k <= 0 {
		return nan()
	}

	if v < 0 {
		p := bsprice(-v, t, x, k, r, q, o)
		i := intrinsic(t, x, k, r, q, o)
		e := p - i
		return i - e
	} else if v*sqrt(t) < Tiny {
		return intrinsic(t, x, k, r, q, o)
	}

	dfq, dfr := exp(-q*t), exp(-r*t)
	d1 := getd1(v, t, x, k, r, q)
	d2 := d1tod2(d1, v, t)
	stdgauss := gaussian.NewGaussian(0, 1)
	Nd1, Nd2 := stdgauss.Cdf(d1), stdgauss.Cdf(d2)

	price = dfq*Nd1*x - dfr*Nd2*k
	if o == Put {
		price -= dfq*x - dfr*k
	} else if o == Straddle {
		price = 2*price - (dfq*x - dfr*k)
	}
	return
}

func bsdelta(v, t, x, k, r, q float64, o OptionType) (delta float64) {

	stdgauss := gaussian.NewGaussian(0, 1)
	Nd1 := stdgauss.Cdf(getd1(v, t, x, k, r, q))

	switch o {
	case Call:
		delta = exp(-q*t) * Nd1
	case Put:
		delta = exp(-q*t) * (Nd1 - 1)
	case Straddle:
		delta = exp(-q*t) * (2*Nd1 - 1)
	default:
		delta = nan()
	}

	return
}

func bsDualDelta(v, t, x, k, r, q float64, o OptionType) (dualdelta float64) {

	stdgauss := gaussian.NewGaussian(0, 1)
	Nd2 := stdgauss.Cdf(getd2(v, t, x, k, r, q))
	switch o {
	case Call:
		dualdelta = -exp(-r*t) * Nd2
	case Put:
		dualdelta = exp(-r*t) * (1 - Nd2)
	case Straddle:
		dualdelta = exp(-q*t) * (1 - 2*Nd2)
	default:
		dualdelta = nan()
	}

	return

}

func bsgamma(v, t, x, k, r, q float64, o OptionType) (gamma float64) {
	d1 := getd1(v, t, x, k, r, q)
	gamma = exp(-q*t-d1*d1/2) / x / v / sqrt(t) * InvSqrt2PI
	if o == Straddle {
		gamma *= 2
	}
	return
}

func bstheta(v, t, x, k, r, q float64, o OptionType) (theta float64) {
	stdgauss := gaussian.NewGaussian(0, 1)
	d1 := getd1(v, t, x, k, r, q)
	d2 := d1tod2(d1, v, t)
	theta = -exp(-q*t) * x * exp(-d1*d1/2) * v / 2 / sqrt(t)
	theta += q*x*exp(-q*t)*stdgauss.Cdf(d1) - r*k*exp(-r*t)*stdgauss.Cdf(d2)
	if o == Straddle {
		theta *= 2
	}
	return
}

func bsvega(v, t, x, k, r, q float64, o OptionType) (vega float64) {
	d1 := getd1(v, t, x, k, r, q)
	vega = x * exp(-q*t-d1*d1/2) * sqrt(t) * InvSqrt2PI
	if o == Straddle {
		vega *= 2
	}
	return
}

func d1tod2(d1, v, t float64) float64 {
	return d1 - v*sqrt(t)
}

func atmApprox(v, t, x, q float64, o OptionType) (price float64) {
	price = exp(-q*t) * v * x * sqrt(t) * InvSqrt2PI
	if o == Straddle {
		price *= 2
	}
	return
}

func fsDeltaApprox(v, t, q float64, o OptionType) float64 {
	delta := exp(-q*t) * v * sqrt(t) * InvSqrt2PI
	if o == Straddle {
		delta *= 2
	}
	return delta
}

func getd1(v, t, x, k, r, q float64) float64 {
	return (log(x/k) + (r-q+0.5*v*v)*t) / v / sqrt(t)
}

func getd2(v, t, x, k, r, q float64) float64 {
	return (log(x/k) + (r-q-0.5*v*v)*t) / v / sqrt(t)
}

func intrinsic(t, x, k, r, q float64, o OptionType) float64 {

	f := exp(-q*t)*x - exp(-r*t)*k
	switch o {
	case Call:
		return max(0, +f)
	case Put:
		return max(0, -f)
	case Straddle:
		return math.Abs(f)
	}
	return nan()
}

func npdf(x float64) float64 {
	return exp(-x*x/2) * InvSqrt2PI
}
