package statistics

// - CumSum function for []float64
// CumSum calculates the cumulative sum of a given slice of float64 values
func CumSum(slice []float64) []float64 {
	sum := 0.0
	cumSum := make([]float64, len(slice))
	for i, val := range slice {
		sum += val
		cumSum[i] = sum
	}
	return cumSum
}

// tryVersion
func TryCumSum(slice []float64) ([]float64, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return CumSum(slice), err
}

// - CumProd function for []float64
// CumProd calculates the cumulative product of a given slice of float64 values
func CumProd(slice []float64) []float64 {
	product := 1.0
	cumProd := make([]float64, len(slice))
	for i, val := range slice {
		product *= val
		cumProd[i] = product
	}
	return cumProd
}

// tryVersion
func TryCumProd(slice []float64) ([]float64, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return CumProd(slice), err
}

