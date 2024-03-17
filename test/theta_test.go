package thetatest

import (
	"math"
	"testing"

	bs "github.com/uscott/go-blackscholes"
)

func Test_Theta(t *testing.T) {

	const N = 15
	eps := 1.0
	var v, tau, x, k, r, q float64 = 1, 1.0 / 12, 100, 110, 0.0, 0.0
	o := bs.Call

	theta := bs.BSTheta(v, tau, x, k, r, q, o) / 365

	t.Logf("\nTheta = %6.2f\n", theta)

	for i := 0; i < N; i++ {

		thetanum := bs.BSThetaNum(v, tau, x, k, r, q, o, eps) / 365

		if math.IsNaN(thetanum) {
			t.Fatal("NaN")
		}

		err := thetanum - theta

		t.Logf(
			"Epsilon = %10.4g, ThetaNum = %6.2f, Error = %8.4f",
			eps, thetanum, err,
		)

		eps /= 2
	}
}
