package implvoltest

import (
	"math/rand"
	"testing"
	"time"

	bs "github.com/uscott/go-blackscholes"
)

func Test_ImpliedVol(t *testing.T) {

	const N int = 50
	tau, x, r, q := 1.0, 100.0, 0.1, 0.05
	var (
		o bs.OptionType
		v float64
	)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < N; i++ {

		k := x + (x/2)*(2*rand.Float64()-1)

		if rand.Float64() < 0.5 {
			v = 0.25 + 1.25*rand.Float64()
		} else {
			v = -0.2 + 0.1*rand.Float64()
		}

		u := rand.Float64()

		switch {
		case u < 0.33:
			o = bs.Call
		case u < 0.67:
			o = bs.Put
		default:
			o = bs.Straddle
		}

		premium := bs.BSPrice(v, tau, x, k, r, q, o)

		pars := &bs.ImpliedVolParams{
			Premium:    premium,
			TimeToExp:  tau,
			Underlying: x,
			Strike:     k,
			Rate:       r,
			Dividend:   q,
			Type:       o,
		}

		implvol, err := bs.ImpliedVol(pars)

		if err != nil {
			t.Fatal(err)
		}

		diff := implvol - v

		t.Logf(
			"Strike = %8.3f, Premium = %8.3f, Vol = %8.4f, ImplVol = %8.4f, Err = %12.6f\n",
			k, premium, v, implvol, diff,
		)

	}
}
