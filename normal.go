package blackscholes

import (
	"fmt"
	"math"
)

func NormCDF(x float64) float64 {
	return 0.5 + 0.5*math.Erf(x/math.Sqrt2)
}

func NormCDFInverse(q float64) (float64, error) {
	if q < 0 || 1 < q {
		return math.NaN(), fmt.Errorf("NormCDFInverse: q must be in the range [0, 1]: %f", q)
	}
	return math.Sqrt2 * math.Erfinv(2*q-1), nil
}
