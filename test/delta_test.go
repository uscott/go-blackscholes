package blackscholes_test

import (
	"math"
	"testing"

	bs "github.com/uscott/go-blackscholes"
)

func Test_Delta(t *testing.T) {

	const N = 15
	eps := 10.0
	var v, tau, x, k, r, q float64 = 0.5, 1, 100, 120, 0.1, 0.5
	o := bs.Call

	delta := bs.BSDelta(v, tau, x, k, r, q, o)

	t.Logf("\nDelta = %6.2f %%\n", 100*delta)

	for i := 0; i < N; i++ {

		deltanum := bs.BSDeltaNum(v, tau, x, k, r, q, o, eps)

		if math.IsNaN(deltanum) {
			t.Fatal("NaN")
		}

		err := deltanum - delta

		t.Logf(
			"Epsilon = %10.4g, DeltaNum = %6.2f %%, Error = %8.4f %%",
			eps, 100*deltanum, 100*err,
		)

		eps /= 2
	}

	x, k = .01, .011

	delta = bs.BSDelta(v, tau, x, k, r, q, o)
	eps = 10

	t.Logf("\nDelta = %6.2f %%\n", 100*delta)

	for i := 0; i < N; i++ {

		deltanum := bs.BSDeltaNum(v, tau, x, k, r, q, o, eps)

		if math.IsNaN(deltanum) {
			t.Fatal("NaN")
		}

		err := deltanum - delta

		t.Logf(
			"Epsilon = %10.4g, DeltaNum = %6.2f %%, Error = %8.4f %%",
			eps, 100*deltanum, 100*err,
		)

		eps /= 2
	}

}
