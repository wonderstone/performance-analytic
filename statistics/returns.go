package statistics

// use gonum package to implement

// - define a struct to calculate the all kinds of returns
type ReturnsCalculator struct {
	R []float64
}

type OptionReturns func(*ReturnsCalculator)

// * for return series
func WithReturns(r []float64) OptionReturns {
	return func(rc *ReturnsCalculator) {
		rc.R = r
	}
}

// * for price or value series
func WithPrices(p []float64) OptionReturns {
	return func(rc *ReturnsCalculator) {
		rc.R = make([]float64, len(p)-1)
		for i := 1; i < len(p); i++ {
			rc.R[i-1] = (p[i] - p[i-1]) / p[i-1]
		}
	}
}

// * for leverage
func WithMultiplier(r []float64, m float64) OptionReturns {
	return func(rc *ReturnsCalculator) {
		rc.R = make([]float64, len(r))
		for i := range r {
			rc.R[i] = r[i] * m
		}
	}
}

func NewReturnsCalculator(opts ...OptionReturns) *ReturnsCalculator {
	rc := &ReturnsCalculator{}
	for _, opt := range opts {
		opt(rc)
	}
	return rc
}

// - Method for excess
// take an interface{} as which can be a float64 or a []float64
// when Rb is a float64 value, it means the risk-free rate after scaling
func (rc *ReturnsCalculator) Excess(Rb interface{}) []float64 {

	result := make([]float64, len(rc.R))
	switch v := Rb.(type) {
	case float64:
		for i := range rc.R {
			result[i] = rc.R[i] - v
		}
		return result
	case []float64:
		for i := range rc.R {
			result[i] = rc.R[i] - v[i]
		}
		return result
	}
	return nil
}

// - Method for cumulative
func (rc *ReturnsCalculator) Cumulative(geometric bool) float64 {
	if geometric {
		// Implement geometric cumulative return calculation
		// iterate over the returns and calculate the product
		res := 1.0
		for _, r := range rc.R {
			res *= 1 + r
		}
		return res - 1
	} else {
		// Implement arithmetic cumulative return calculation
		res := 0.0
		for _, r := range rc.R {
			res += r
		}
		return res
	}
}
