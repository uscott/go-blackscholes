package blackscholes

import "math"

func NormCDF(x float64) float64 {
	return 0.5 + 0.5*math.Erf(x/math.Sqrt2)
}

func NormCDFInverse(q float64) float64 {
	return math.Sqrt2 * math.Erfinv(2*q-1)
}
