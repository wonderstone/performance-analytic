package statistics

import (
	"math"

	"github.com/gonum/floats"
	"gonum.org/v1/gonum/mat"
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

// - ActivePremium function
func ActivePremium(Ra []float64, Rb []float64, scale int, geometric bool) float64 {
	ARa := AnnualizedReturn(Ra, scale, geometric)
	ARb := AnnualizedReturn(Rb, scale, geometric)
	return ARa - ARb
}

// - StdDev function
func StdDev(data []float64) float64 {
	return stat.StdDev(data, nil)
}

// - StdDevAnnualized function
func StdDevAnnualized(data []float64, scale int) float64 {
	return StdDev(data) * math.Sqrt(float64(scale))
}

// - SharpeRatio function
func SharpeRatio(Ra []float64, Rf interface{}, scale int, geometric bool) float64 {
	rts := ReturnsCalculator{Ra}

	// calculate the excess returns
	ExcessRa := rts.Excess(Rf)

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

// - Drawdowns function
// todo: seems useless at this moment
func Drawdowns(Ra []float64) []float64 {
	cumulative := CumProdAdd(Ra, 1.0)
	// add the initial value of 1.0
	maxcumulative := CumMax(append([]float64{1.0}, cumulative...))[1:]
	// vectorize the drawdown calculation
	drawdowns := make([]float64, len(Ra))
	for i := range Ra {
		drawdowns[i] = cumulative[i]/maxcumulative[i] - 1
	}
	return drawdowns
}

// - TrackingError function
func TrackingError(Ra, Rb []float64, scale int) float64 {
	// calculate the excess returns
	rt := ReturnsCalculator{Ra}
	ExcessRa := rt.Excess(Rb)

	// calculate the standard deviation of the excess returns
	return stat.StdDev(ExcessRa, nil) * math.Sqrt(float64(scale))
}

// - InformationRatio function
// This relates the degree to which an investment has beaten
// the benchmark to the consistency with
// which the investment has beaten the benchmark.
func InformationRatio(Ra, Rb []float64, scale int) float64 {
	// Active premium
	ap := ActivePremium(Ra, Rb, scale, true)
	// tracking error
	te := TrackingError(Ra, Rb, scale)
	// calculate the Information Ratio
	return ap / te
}

// - DownsideDeviation function
func DownsideDeviation(data []float64, threshold float64, tag string) float64 {
	// calculate the downside returns
	downside := make([]float64, 0)
	for _, r := range data {
		if r < threshold {
			downside = append(downside, r-threshold)
		}
	}
	var length float64
	// calculate the downside deviation
	if tag == "subset" {
		length = float64(len(downside))
	} else {
		length = float64(len(data))
	}

	sum := 0.0
	for _, value := range downside {
		sum += math.Pow(value, 2)
	}

	return math.Sqrt(sum / length)
}

// - DownsideVariance function
func DownsideVariance(data []float64, threshold float64, tag string) float64 {
	// calculate the downside returns
	downside := make([]float64, 0)
	for _, r := range data {
		if r < threshold {
			downside = append(downside, r-threshold)
		}
	}
	var length float64
	// calculate the downside deviation
	if tag == "subset" {
		length = float64(len(downside))
	} else {
		length = float64(len(data))
	}

	sum := 0.0
	for _, value := range downside {
		sum += math.Pow(value, 2)
	}

	return sum / length
}

// - DownsidePotential function
func DownsidePotential(data []float64, threshold float64, tag string) float64 {
	// calculate the downside returns
	downside := make([]float64, 0)
	for _, r := range data {
		if r < threshold {
			downside = append(downside, r-threshold)
		}
	}
	var length float64
	// calculate the downside deviation
	if tag == "subset" {
		length = float64(len(downside))
	} else {
		length = float64(len(data))
	}

	sum := 0.0
	for _, value := range downside {
		sum += value
	}

	return sum / length
}

// - Hurst index function
// A Hurst index between 0.5 and 1 suggests that the returns are persistent. At 0.5, the index suggests returns are totally
// random. Between 0 and 0.5 it suggests that the returns are mean reverting.
func HurstIndex(data []float64) float64 {
	n := float64(len(data))

    min, max := floats.Min(data), floats.Max(data)
    sd := stat.StdDev(data, nil)

    m := (max - min) / sd
    result := math.Log(m) / math.Log(n)

    return result
}

// - MarketTiming function
// TH for Treynor-Mazuy model
// HM for Henriksson-Merton model
func MarketTiming(Ra,Rb []float64,Rf interface{}, tag string) (Alpha, Beta ,Gamma float64) {
	rtsa := ReturnsCalculator{Ra}
	rtsb := ReturnsCalculator{Rb}

	// calculate the excess returns
	ExcessRa := rtsa.Excess(Rf)
	ExcessRb := rtsb.Excess(Rf)
	// make a slice with the same length as the returns
	D := make([]float64, len(ExcessRb))
	switch tag {
	case "TH":
		copy(D,ExcessRb)
	case "HM":
		for i := range ExcessRb {
			if ExcessRb[i] < 0 {
				D[i] = 1
			} else {
				D[i] = 0
			}
		}
	default:
		copy(D,ExcessRb)
	}
	// let ExcessRa is y 
	
	Y := mat.NewDense(len(ExcessRa), 1, ExcessRa)
	// let ExcessRb and ExcessRb*D as Xs
	Xs := mat.NewDense(len(ExcessRb), 3, nil)
	for i := 0; i < len(ExcessRb); i++ {
		Xs.Set(i, 0, 1)
		Xs.Set(i, 1, ExcessRb[i])
		Xs.Set(i, 2, ExcessRb[i]*D[i])
	}
	// calculate the coefficients
	olsres := NewOLS(Xs,Y)

	olsres.Run()

	cos := olsres.Coefficients().RawMatrix().Data
	Alpha = cos[0]
	Beta = cos[1]
	Gamma = cos[2]

	return 
}