package blackscholes

func BSDeltaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	e := abs(eps)
	pu := BSPriceNoErrorCheck(v, t, x+e, k, r, q, o)

	if x < e {
		pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)
		pd := ZeroUnderlyingBSPrice(t, k, r, o)
		cu, cm, cd := x/e/(e+x), (e-x)/e/x, -e/x/(e+x)
		return cu*pu + cm*pm + cd*pd
	}

	pd := BSPriceNoErrorCheck(v, t, x-e, k, r, q, o)

	return (pu - pd) / 2 / e
}

func BSGammaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	e := abs(eps)

	if x < e {
		puu := BSPriceNoErrorCheck(v, t, x+e, k, r, q, o)
		pu := BSPriceNoErrorCheck(v, t, 2*x, k, r, q, o)
		pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)
		pd := ZeroUnderlyingBSPrice(t, k, r, o)
		d := (x + e) * (x + e)
		cuu, cu, cm, cd := 2*x/e/d, (e-x)/x/d, -2/e/x, 1/x/(e+x)
		return cuu*puu + cu*pu + cm*pm + cd*pd
	}
	pu := BSPriceNoErrorCheck(v, t, x+e, k, r, q, o)
	pd := BSPriceNoErrorCheck(v, t, x-e, k, r, q, o)
	pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)

	return (pu - 2*pm + pd) / (e * e)
}

func BSVegaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	pu := BSPriceNoErrorCheck(v+eps, t, x, k, r, q, o)
	pd := BSPriceNoErrorCheck(v-eps, t, x, k, r, q, o)

	return (pu - pd) / 2 / eps
}

func BSThetaNum(v, t, x, k, r, q float64, o OptionType, eps float64) float64 {

	if CheckPriceParams(t, x, k, o) != nil {
		return nan()
	}

	e := abs(eps)

	if t < e {
		pu := Intrinsic(0, x, k, 0, 0, o)
		pd := BSPriceNoErrorCheck(v, t+e, x, k, r, q, o)
		pm := BSPriceNoErrorCheck(v, t, x, k, r, q, o)
		wu, wd := e/(e+t), t/(e+t)
		return wu*(pu-pm)/(e+t) + wd*(pm-pd)/(e+t)
	}

	pu := BSPriceNoErrorCheck(v, t-e, x, k, r, q, o)
	pd := BSPriceNoErrorCheck(v, t+e, x, k, r, q, o)

	return (pu - pd) / 2 / e
}
