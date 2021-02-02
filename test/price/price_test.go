package pricetest

import (
	"math"
	"testing"

	"github.com/uscott/go-blackscholes"
)

func Test_Price(t *testing.T) {

	v, tt, x, k, r, q := 0.5, 1.0/12, 100.0, 120.0, 0.1, 0.05
	o := blackscholes.Call

	N := 20
	price := blackscholes.BSPrice(v, tt, x, k, r, q, o)

	t.Logf("Price = %.2f\n", price)

	var (
		nprev   uint
		simprev float64
	)

	t.Log("Sim:")
	for i := 1; i <= N; i++ {

		n := uint(1 << i)
		simprice := blackscholes.BSPriceSim(v, tt, x, k, r, q, o, n)
		w, wprev := float64(n)/float64(n+nprev), float64(nprev)/float64(n+nprev)
		simprice = w*simprice + wprev*simprev
		err := simprice/price - 1
		nprev += n
		simprev = simprice

		t.Logf(
			"Number of randoms = %8d,\tsim price = %6.2f,\t|error| = %8.3f %%\n",
			nprev, simprice, 100*math.Abs(err))
	}

}
