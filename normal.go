package blackscholes

import (
	"gonum.org/v1/gonum/stat/distuv"
)

var normal = distuv.Normal{Mu: 0, Sigma: 1}

func NormCDF(x float64) float64 {
	return normal.CDF(x)
}

func NormPDF(x float64) float64 {
	return normal.Prob(x)
}

func NormCDFInverse(q float64) float64 {
	return normal.Quantile(q)
}
