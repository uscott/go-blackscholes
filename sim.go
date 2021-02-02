package blackscholes

import (
	"math/rand"
	"sync"
	"time"
)

func BSPriceSim(v, t, x, k, r, q float64, o OptionType, n uint) float64 {

	if !ValidOptionType(o) || n == 0 {
		return nan()
	}

	rand.Seed(time.Now().UnixNano())

	mu, wg := new(sync.Mutex), new(sync.WaitGroup)
	sum, x0 := 0.0, exp(-q*t)*x
	m, s := x0*exp((r-0.5*v*v)*t), v*sqrt(t)

	wg.Add(int(n))
	for i := 0; i < int(n); i++ {
		go func(i int) {
			mu.Lock()
			e := exp(s * rand.NormFloat64())
			x = m * e
			sum += Intrinsic(0, x, k, 0, 0, o)
			x = m / e
			sum += Intrinsic(0, x, k, 0, 0, o)
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	return exp(-r*t) * sum / float64(2*n)
}
