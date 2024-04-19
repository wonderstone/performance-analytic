package statistics

import "gonum.org/v1/gonum/stat"

// use gonum package to implement

// - define a struct to calculate the CAPM

type CAPM struct {
	Ra []float64
	Rb []float64
}

type OptionCAPM func(*CAPM)

// * for portfolio return series
func WithRa(ra []float64) OptionCAPM {
	return func(c *CAPM) {
		c.Ra = ra
	}
}

// * for benchmark return series
func WithRb(rb []float64) OptionCAPM {
	return func(c *CAPM) {
		c.Rb = rb
	}
}

func NewCAPM(opts ...OptionCAPM) *CAPM {
	c := &CAPM{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}	

// - Method for Beta
func (c *CAPM) Beta() (float64) {
	return CoVariance(c.Ra, c.Rb) / Variance(c.Rb)
}

// - Method for Alpha
func (c *CAPM) Alpha(rf interface{}) float64 {
	beta := c.Beta()
	ra := ReturnsCalculator{c.Ra}
	rb := ReturnsCalculator{c.Rb}	
	excessRa := ra.Excess(rf)
	excessRb := rb.Excess(rf)
	return stat.Mean(excessRa,nil) - beta *stat.Mean(excessRb,nil)

}



// - Method for TimingRatio
func (c *CAPM) TimingRatio() float64 {
	// sort out the positive and negative returns into different slices 
	positiveRa := make([]float64, 0)
	negativeRa := make([]float64, 0)
	positiveRb := make([]float64, 0)
	negativeRb := make([]float64, 0)
	for i, val := range c.Rb {
		if val > 0 {
			positiveRb = append(positiveRb, val)
			positiveRa = append(positiveRa, c.Ra[i])
		} else {
			negativeRb = append(negativeRb, val)
			negativeRa = append(negativeRa, c.Ra[i])
		}
	}
	// give betas to the positive and negative returns
	betaPositive := CoVariance(positiveRa, positiveRb) / Variance(positiveRb)
	betaNegative := CoVariance(negativeRa, negativeRb) / Variance(negativeRb)
	// calculate the timing ratio
	return betaPositive/betaNegative	
}

