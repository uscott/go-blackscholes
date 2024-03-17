package blackscholes_test

import (
	"math"
	"testing"

	bs "github.com/uscott/go-blackscholes"
)

func Test_Gamma(t *testing.T) {

	const N = 15
	eps := 10.0
	var v, tau, x, k, r, q float64 = 0.5, 1, 100, 80, 0.1, 0.5
	o := bs.Put

	gamma := bs.BSGamma(v, tau, x, k, r, q, o)

	t.Logf("\n$ Gamma = %6.2f\n", x*gamma)

	for i := 0; i < N; i++ {

		gammanum := bs.BSGammaNum(v, tau, x, k, r, q, o, eps)

		if math.IsNaN(gammanum) {
			t.Fatal("NaN")
		}

		err := x * (gammanum - gamma)

		t.Logf(
			"Epsilon = %10.4g, $ GammaNum = %6.2f, Error = %8.4f",
			eps, x*gammanum, err,
		)

		eps /= 2
	}

	x, k = .1, .09

	gamma = bs.BSGamma(v, tau, x, k, r, q, o)
	eps = 10

	t.Logf("\n$ Gamma = %6.2f\n", x*gamma)

	for i := 0; i < N; i++ {

		gammanum := bs.BSGammaNum(v, tau, x, k, r, q, o, eps)

		if math.IsNaN(gammanum) {
			t.Fatal("NaN")
		}

		err := x * (gammanum - gamma)

		t.Logf(
			"Epsilon = %10.4g, GammaNum = %6.2f, Error = %8.4f",
			eps, x*gammanum, err,
		)

		eps /= 2
	}
}
