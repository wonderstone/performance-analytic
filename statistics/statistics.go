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
	return ActivePremium(Ra, Rb, scale, geometric), err
}

// - StdDev function
func StdDev(data []float64) float64 {
	return stat.StdDev(data, nil)
}

// TryVersion
func TryStdDev(data []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return StdDev(data), err
}

// - StdDevAnnualized function
func StdDevAnnualized(data []float64, scale int) float64 {
	return StdDev(data) * math.Sqrt(float64(scale))
}
// TryVersion
func TryStdDevAnnualized(data []float64, scale int) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return StdDevAnnualized(data, scale), err
}



// - SharpeRatio function
func SharpeRatio(Ra []float64, Rf interface{}, scale int, geometric bool) float64 {
	rts := ReturnsCalculator{Ra}
	
	// calculate the excess returns
	ExcessRa:= rts.Excess(Rf)

	// calculate the mean excess return
	MeanExcessReturn := stat.Mean(ExcessRa, nil)
	

	// calculate the ExcessRa standard deviation
	// ! performance analytics package is wrong on this calculation
	// ! the correct calculation is the standard deviation of the excess returns
	// ! although when Rf has no volatility, the standard deviations are the same
	// ! https://en.wikipedia.org/wiki/Sharpe_ratio
	// ! https://www.investopedia.com/terms/s/sharperatio.asp
	// ! or even the PerformanceAnalytics package notes @ page 173
	// ! https://cran.r-project.org/web/packages/PerformanceAnalytics/PerformanceAnalytics.pdf
	// StdDevRa := stat.StdDev(Ra, nil)
	StdDevRa := stat.StdDev(ExcessRa, nil)
	
	// calculate the Sharpe Ratio
	return MeanExcessReturn / StdDevRa
}

// TryVersion
func TrySharpeRatio(Ra []float64, Rf interface{}, scale int, geometric bool) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return SharpeRatio(Ra, Rf, scale, geometric), err
}

// - MaxDrawdown function
func MaxDrawdown(Ra []float64) float64 {
	// calculate the cumulative returns
	cumulative := 1.0
	tmpMaxCumulative := 0.0
	maxDrawdown := 0.0
	for _, r := range Ra {
		cumulative *= 1 + r
		if cumulative > tmpMaxCumulative {
			tmpMaxCumulative = cumulative
		}
		drawdown := 1 - cumulative/tmpMaxCumulative
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}
	return maxDrawdown
}

// TryVersion
func TryMaxDrawdown(Ra []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return MaxDrawdown(Ra), err
}

// 