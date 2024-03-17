package blackscholes

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

const defaultNumPaths uint = 1 << 16

func PriceSim(vol, timeToExpiry, spot, strike, interestRate, dividendYield float64, optionType OptionType, numPaths ...uint) (price float64, err error) {

	npaths := defaultNumPaths
	if len(numPaths) > 0 {
		npaths = numPaths[0]
	}

	price = math.NaN()

	if npaths == 0 {
		err = errors.New("number of paths must be positive")
		return
	}

	if !ValidOptionType(optionType) {
		err = ErrUnknownOptionType
		return
	}

	sum := 0.0
	expectedSpot := spot * math.Exp((interestRate-dividendYield)*timeToExpiry)
	sigma := vol * math.Sqrt(timeToExpiry)
	mu := -0.5 * sigma * sigma

	mtx, wg := new(sync.Mutex), new(sync.WaitGroup)
	wg.Add(int(npaths - 1))

	for i := uint(0); i < npaths-1; i += 2 {

		go func() {

			defer mtx.Unlock()
			mtx.Lock()

			defer wg.Done()

			z := rand.NormFloat64()

			spot = expectedSpot * math.Exp(mu+sigma*z)
			sum += Intrinsic(0, spot, strike, 0, 0, optionType)

			spot = expectedSpot * math.Exp(mu-sigma*z)
			sum += Intrinsic(0, spot, strike, 0, 0, optionType)

		}()
	}

	wg.Wait()

	if npaths%2 == 1 {
		z := rand.NormFloat64()
		spot = expectedSpot * math.Exp(mu+sigma*z)
		sum += Intrinsic(0, spot, strike, 0, 0, optionType)
	}

	price = math.Exp(-interestRate*timeToExpiry) * sum / float64(npaths)

	return
}
