package blackscholes

import (
	"fmt"
	"math"
	"time"
)

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

func ImpliedVol(pars *ImpliedVolParams) (vol float64, err error) {

	if pars == nil {
		return nan(), ErrNilPtrArg
	}

	p, t, x, k, r, q := GetFloatVolParams(pars)
	o := pars.Type
	if err = CheckPriceParams(t, x, k, o); err != nil {
		return nan(), err
	}

	if t == 0 || x == 0 || k == 0 {
		return 0, nil
	}

	intrval := Intrinsic(t, x, k, r, q, o)
	extrval := p - intrval

	if abs(extrval) <= math.SmallestNonzeroFloat64 {
		return 0, nil
	}

	lb, ub, tol, maxit := GetVolSearchParams(pars)

	CheckVolSearchParams(&lb, &ub, &tol, &maxit)

	var (
		it             int
		plo, phi, pmid float64
	)
	for it = 0; it < maxit; it++ {
		plo = BSPriceNoErrorCheck(lb, t, x, k, r, q, o)
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
		phi = BSPriceNoErrorCheck(ub, t, x, k, r, q, o)
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

	for it = 0; it < maxit; it++ {

		vol = 0.5 * (lb + ub)
		pmid = BSPriceNoErrorCheck(vol, t, x, k, r, q, o)

		switch {
		case ub-lb < tol, pmid == p:
			CorrectVolSign(pmid-intrval, &vol)
			return
		case p < pmid:
			ub = vol
		case pmid < p:
			lb = vol
		}
	}

	plo, phi = BSPrice(lb, t, x, k, r, q, o), BSPrice(ub, t, x, k, r, q, o)
	return nan(), fmt.Errorf(
		"Did not converge - lb, ub, lb price, ub price, mid, iters: %v, %v, %v, %v, %v, %d",
		lb, ub, plo, phi, pmid, it,
	)
}

func CheckVolSearchParams(lb, ub, tol *float64, maxit *int) {

	if lb == nil || ub == nil || tol == nil || maxit == nil {
		panic(ErrNilPtrArg)
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

func CorrectVolSign(extrinsic float64, vol *float64) {
	if vol == nil {
		panic(ErrNilPtrArg)
	}
	if extrinsic > 0 && *vol < 0 || extrinsic < 0 && *vol > 0 {
		*vol = -*vol
		return
	}
	if extrinsic == 0 {
		*vol = 0
	}
}

func GetFloatVolParams(pars *ImpliedVolParams) (p, t, x, k, r, q float64) {
	if pars == nil {
		panic(ErrNilPtrArg)
	}
	t = pars.Expiry.Sub(pars.Now).Hours() / 24 / 365
	p, x, k, r, q = pars.Premium,
		pars.Underlying, pars.Strike, pars.Rate, pars.Dividend
	return
}

func GetVolSearchParams(pars *ImpliedVolParams) (lb, ub, tol float64, maxit int) {
	if pars == nil {
		panic(ErrNilPtrArg)
	}
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
