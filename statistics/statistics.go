package statistics


import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// - Annualized Return function
// scale number of periods in a year
// daily scale = 252, weekly scale = 52,
// monthly scale = 12, quarterly scale = 4
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
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return AnnualizedReturn(Ra, scale, geometric), err
}



// - ActivePremium function
func ActivePremium(Ra []float64, Rb []float64, scale int, geometric bool) float64 {
	ARa := AnnualizedReturn(Ra, scale, geometric)
	ARb := AnnualizedReturn(Rb, scale, geometric)
	return ARa - ARb
}

func TryActivePremium(Ra, Rb []float64, scale int, geometric bool) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return ActivePremium(Ra, Rb,scale,geometric), err
}

// - SharpeRatio function
func SharpeRatio(returns []float64, riskFreeRate float64, scale int, geometric bool) float64 {
	// Implement Annualized Return for the returns
	AR := AnnualizedReturn(returns, scale, geometric)
	// Calculate the standard deviation of the returns
	stdDev := stat.StdDev(returns, nil)
	// Implement Sharpe Ratio calculation
	return (AR - riskFreeRate) / stdDev
}