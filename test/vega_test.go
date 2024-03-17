package blackscholes_test

import (
	"math"
	"testing"

	bs "github.com/uscott/go-blackscholes"
)

func Test_Vega(t *testing.T) {

	const N = 15
	eps := 10.0
	var v, tau, x, k, r, q float64 = 0.5, 1, 100, 120, 0.1, 0.5
	o := bs.Call

	vega := bs.BSVega(v, tau, x, k, r, q, o) / 100

	t.Logf("\nVega = %7.3f\n", vega)

	for i := 0; i < N; i++ {

		veganum := bs.BSVegaNum(v, tau, x, k, r, q, o, eps) / 100

		if math.IsNaN(veganum) {
			t.Fatal("NaN")
		}

		err := veganum - vega

		t.Logf(
			"Epsilon = %10.4g, VegaNum = %7.3f, Error = %8.4f",
			eps, veganum, err,
		)

		eps /= 2
	}

}
