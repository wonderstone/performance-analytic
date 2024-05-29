package statistics

import (
	"encoding/csv"
	"errors"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	// stat2 "github.com/grd/stat"
)

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

// - CumProdAdd function for []float64
func CumProdAdd(slice []float64, other interface{}) []float64 {
	switch v := other.(type) {
	case float64:
		product := 1.0
		cumProd := make([]float64, len(slice))
		for i, val := range slice {
			product *= val + v
			cumProd[i] = product
		}
		return cumProd
	case []float64:
		product := 1.0
		cumProd := make([]float64, len(slice))
		for i, val := range slice {
			product *= val + v[i]
			cumProd[i] = product
		}
		return cumProd
	default:
		panic(errors.New("invalid type"))
	}
}

// - CumMax function for []float64
// CumMax calculates the cumulative maximum of a given slice of float64 values
func CumMax(slice []float64) []float64 {
	max := math.Inf(-1)
	cumMax := make([]float64, len(slice))
	for i, val := range slice {
		if val > max {
			max = val
		}
		cumMax[i] = max
	}
	return cumMax
}

// - CumMin function for []float64
// CumMin calculates the cumulative minimum of a given slice of float64 values
func CumMin(slice []float64) []float64 {
	min := math.Inf(1)
	cumMin := make([]float64, len(slice))
	for i, val := range slice {
		if val < min {
			min = val
		}
		cumMin[i] = min
	}
	return cumMin
}

// - Variance function
// Variance calculates the variance of a given slice of float64 values
func Variance(data []float64) float64 {
	return stat.Variance(data, nil)
}

// - SemiVariance function
func SemiVariance(data []float64, tag string) float64 {
	// calculate the mean of the data
	mean := stat.Mean(data, nil)
	// calculate the semi variance
	return DownsideVariance(data, mean, tag)
}

// - SemiDeviation function
func SemiDeviation(data []float64, tag string) float64 {
	// calculate the mean of the data
	mean := stat.Mean(data, nil)
	// calculate the semi deviation
	return DownsideDeviation(data, mean, tag)
}


// - CoVariance function
// CoVariance calculates the covariance of two given slices of float64 values
func CoVariance(x, y []float64) float64 {
	return stat.Covariance(x, y, nil)
}

// - Correlation function
// Correlation calculates the correlation coefficient between two given slices of float64 values
func Correlation(x, y []float64) float64 {
	return stat.Correlation(x, y, nil)
}

// - Skewness function
// Skewness calculates the skewness of a given slice of float64 values
func Skewness(data []float64, tag string) float64 {
	// has moment、fisher、sample for now
	mean := stat.Mean(data, nil)
	variance := stat.Variance(data, nil)
	n := float64(len(data))

	sum := 0.0
	for _, value := range data {
		sum += math.Pow((value - mean), 3)
	}

	switch tag {
	case "moment":
		return sum / (math.Pow(math.Sqrt(variance*(n-1)/n), 3) * n)
	case "fisher":
		if n < 3 {
			return math.NaN()
		} else {
			sumCube := 0.0
			sumSquare := 0.0
			for _, value := range data {
				sumCube += math.Pow(value, 3)
				sumSquare += math.Pow(value, 2)
			}
			return (math.Sqrt(n*(n-1)) / (n - 2)) * (sumCube / n) / math.Pow((sumSquare/n), 1.5)
		}
	case "sample":
		return sum / math.Pow(math.Sqrt(variance*(n-1)/n), 3) * n / ((n - 1) * (n - 2))

	default:
		return sum / (math.Pow(math.Sqrt(variance*(n-1)/n), 3) * n)
	}
}

// - Kurtosis function
// Kurtosis calculates the kurtosis of a given slice of float64 values

func Kurtosis(x []float64, tag string) float64 {
    n := float64(len(x))
    mean := stat.Mean(x, nil)
    variance := stat.Variance(x, nil)

    sumFourthPower := 0.0
    sumSquare := 0.0
    for _, value := range x {
        sumFourthPower += math.Pow((value - mean), 4)
        sumSquare += math.Pow(value, 2)
    }

    switch tag {
    case "excess":
        return sumFourthPower/(math.Pow(variance*(n-1)/n, 2)*n) - 3
    case "moment":
        return sumFourthPower / (math.Pow(variance*(n-1)/n, 2) * n)
    case "fisher":
		sumFourthPowerFisher := 0.0
		sumSquareFisher := 0.0
		for _, value := range x {
			sumFourthPowerFisher += math.Pow(value, 4)
			sumSquareFisher += math.Pow(value, 2)
		}
        return ((n+1)*(n-1)*((sumFourthPowerFisher/n)/math.Pow(sumSquareFisher/n, 2)-(3*(n-1))/(n+1)))/((n-2)*(n-3))
    case "sample":
        return sumFourthPower/math.Pow(variance, 2) * n * (n+1) / ((n-1)*(n-2)*(n-3))
    case "sample_excess":
        return sumFourthPower/math.Pow(variance, 2)*n*(n+1)/((n-1)*(n-2)*(n-3)) - 3*math.Pow(n-1, 2)/((n-2)*(n-3))
    default:
        return sumFourthPower / (math.Pow(variance*(n-1)/n, 2) * n)
    }
}


// * function to read managers.csv data
func ReadData(path string) (dt [][]string, fields []string) {
	// read the data from the csv file using io package
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// get the fields from the csv file
	fields, err = reader.Read()
	if err != nil {
		panic(err)
	}

	reader.FieldsPerRecord = 0

	dt, err = reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return dt, fields
}

// * function to check the position of a field in the fields slice
func CheckPos(fields []string, name string) (pos int, e error) {
	for i, field := range fields {
		if field == name {
			return i, nil
		}
	}
	err := errors.New("field not found")
	return -1, err

}

// * function to get a [][]string 2D slice the second column
// GetSecondDimensionData returns a [][]string slice of the second column of a given 2D slice
func GetSecondDimensionData(data [][]string, index int) []string {
	secondColumn := make([]string, len(data))
	for i, row := range data {
		secondColumn[i] = row[index]
	}
	return secondColumn
}

// * function to change []string to []float64
// TryStringToFloatSlice converts a []string slice to a []float64 slice
func TryStringToFloatSlice(s []string) ([]float64, error) {
	var f []float64
	for _, str := range s {
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		f = append(f, float)
	}
	return f, nil
}

// StringToFloatSlice forcely converts a []string slice to a []float64 slice and returns the result without error handling
func StringToFloatSlice(s []string) []float64 {
	var f []float64
	for _, str := range s {
		float, e := strconv.ParseFloat(str, 64)
		if e != nil {
		} else {
			f = append(f, float)
		}
	}
	return f
}

// StringToFloatSlice forcely converts a []string slice to a []float64 slice with same length for benchmark
func StringToFloatSliceBench(Ra, Rb []string) (RaF, RbF []float64) {
	for i, s := range Ra {
		fa, e := strconv.ParseFloat(s, 64)
		if e != nil {
		} else {
			RaF = append(RaF, fa)
			// assume Rb has all reasonable values
			fb, _ := strconv.ParseFloat(Rb[i], 64)
			RbF = append(RbF, fb)
		}

	}
	return RaF, RbF
}
