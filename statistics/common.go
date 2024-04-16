package statistics

import (
	"encoding/csv"
	"errors"
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


// - Variance function
// Variance calculates the variance of a given slice of float64 values
func Variance(data []float64) float64 {
	return stat.Variance(data, nil)
}

// tryVersion
func TryVariance(data []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return Variance(data), err
}



// - CoVariance function
// CoVariance calculates the covariance of two given slices of float64 values
func CoVariance(x, y []float64) float64 {
	return stat.Covariance(x, y, nil)
}

// tryVersion
func TryCoVariance(x, y []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return CoVariance(x, y), err
}

// - Correlation function
// Correlation calculates the correlation coefficient between two given slices of float64 values
func Correlation(x, y []float64) float64 {
	return stat.Correlation(x, y, nil)
}

// tryVersion
func TryCorrelation(x, y []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
		}()
		return Correlation(x, y), err
}

// - Skewness function
// Skewness calculates the skewness of a given slice of float64 values
func Skewness(data []float64) float64 {
	return stat.Skew(data, nil)
}
// tryVersion
func TrySkewness(data []float64) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return Skewness(data), err
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



