package blackscholes

import (
	"math/rand"
	"sync"
)

func BSPriceSim(v, t, x, k, r, q float64, o OptionType, n uint) float64 {

	if !ValidOptionType(o) || n == 0 {
		return nan()
	}

	mu, wg := new(sync.Mutex), new(sync.WaitGroup)
	sum, x0 := 0.0, exp(-q*t)*x
	m, s := x0*exp((r-0.5*v*v)*t), v*sqrt(t)

	wg.Add(int(n - 1))
	for i := 1; i < int(n); i++ {
		go func(i int) {
			mu.Lock()
			u := (float64(i) - 0.5 + rand.Float64()) / float64(n)
			e := exp(s * NormCDFInverse(u))
			x = m * e
			sum += Intrinsic(0, x, k, 0, 0, o)
			x = m / e
			sum += Intrinsic(0, x, k, 0, 0, o)
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	u := 0.5 * rand.Float64() / float64(n)
	e := exp(s * NormCDFInverse(u))
	x = m * e
	sum += Intrinsic(0, x, k, 0, 0, o)
	x = m / e
	sum += Intrinsic(0, x, k, 0, 0, o)

	return exp(-r*t) * sum / float64(2*n)
}
