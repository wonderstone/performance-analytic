package statistics

// use gonum package to implement:
// def active_premium(returns, benchmark):
// """
// calculate the active premium
// :param returns: the returns
// :param benchmark: the benchmark returns
// :return: the active premium
// """
// diff = returns - benchmark
// return np.mean(diff)

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// - Annualized Return function
// scale number of periods in a year
// daily scale = 252, monthly scale = 12, quarterly scale = 4
func AnnualizedReturn(Ra []float64, scale int, geometric bool) float64 {
	if geometric {
		// Implement geometric annualized return calculation
		// iterate over the returns and calculate the product
		product := 1.0
		for _, r := range Ra {
			product *= 1 + r
		}
		return math.Pow(product, float64(scale)/float64(len(Ra))) - 1
	} else {
		// Implement arithmetic annualized return calculation
		return stat.Mean(Ra, nil) * float64(scale)
	}
}

// try to calculate the annualized return
func TryAnnualizedReturn(Ra []float64, scale int, geometric bool) (res float64, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return AnnualizedReturn(Ra, scale, geometric), err
}

// - ActivePremium function
func ActivePremium(Ra []float64, Rb []float64, scale int, geometric bool) float64 {
	// Implement Annualized Return for both Ra and Rb
	ARa := AnnualizedReturn(Ra, scale, geometric)
	ARb := AnnualizedReturn(Rb, scale, geometric)
	// Implement Active Premium calculation
	return ARa - ARb
}

func TryActivePremium(Ra, Rb []float64, scale int, geometric bool) (res float64, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return ActivePremium(Ra, Rb,scale,geometric), err
}
