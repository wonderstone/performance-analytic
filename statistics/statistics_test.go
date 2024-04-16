package statistics

// test the ActivePremium function using
// testify package
import (
	"encoding/csv"
	"errors"
	"os"
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

// * function to change []string to []float64
var dt [][]string
var fds []string

// ~ init this test
func init() {
	dt, fds = ReadData("../data/managers.csv")
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






// Test ReturnsCalculator all methods
func TestReturnsCalculator(t *testing.T) {
	// define the returns for HAM1
	rtp, _ := CheckPos(fds, "HAM1")
	rts := GetSecondDimensionData(dt, rtp)
	rt, e := TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}
	rc := ReturnsCalculator{rt}

	// test the Excess method
	excessReturns,_ := rc.Excess(0.04/12)
	//  only check the first 7 decimal places
	assert.InDelta(t,excessReturns[0], 0.004066667, 0.0000001)
	assert.Equal(t, len(excessReturns), len(rt))

	// thest the Cumulative method
	CumReturn,_ := rc.Cumulative(true)
	// only check the first 7 decimal places
	assert.InDelta(t,CumReturn, 3.126671, 0.000001)

}



// TestAnnualizedReturn tests the AnnualizedReturn function
func TestAnnualizedReturn(t *testing.T) {

	// define the returns for HAM1
	rtp, _ := CheckPos(fds, "HAM1")
	rts := GetSecondDimensionData(dt, rtp)
	rt, e := TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}

	// calculate the annualized return @ geometric = true
	// only check the first 7 decimal places
	// this number（0.1375320）is from the R code
	AR := AnnualizedReturn(rt, 12, true)
	assert.InDelta(t, AR, 0.1375320, 0.0000001)
	// calculate the annualized return @ geometric = false
	// this number（0.1334727）is from the R code
	AR = AnnualizedReturn(rt, 12, false)
	assert.InDelta(t, AR, 0.1334727, 0.0000001)

	// define the returns for HAM2
	rtp1, _ := CheckPos(fds, "HAM2")
	rts1 := GetSecondDimensionData(dt, rtp1)
	rt1 := StringToFloatSlice(rts1)
	// this number（0.1746569）is from the R code
	AR = AnnualizedReturn(rt1, 12, true)
	assert.InDelta(t, AR, 0.1746569, 0.0000001)
}

// TestActivePremium tests the ActivePremium function
func TestActivePremium(t *testing.T) {
	// define the returns
	rtp, _ := CheckPos(fds, "HAM1")
	rts := GetSecondDimensionData(dt, rtp)
	rt, e := TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}
	// define the benchmark returns
	bmp, _ := CheckPos(fds, "SP500 TR")
	bms := GetSecondDimensionData(dt, bmp)
	bm, e := TryStringToFloatSlice(bms)
	if e != nil {
		panic(e)
	}
	// calculate the active premium
	// this number（0.04078668）is from the R code
	AP := ActivePremium(rt, bm, 12, true)
	assert.InDelta(t, AP, 0.04078668, 0.0000001)

	// define the returns for HAM2
	rtp1, _ := CheckPos(fds, "HAM2")
	rts1 := GetSecondDimensionData(dt, rtp1)
	rt1,bm1 := StringToFloatSliceBench(rts1,bms)

	// calculate the active premium
	// this number（0.07759873）is from the R code
	AP = ActivePremium(rt1, bm1, 12, true)
	assert.InDelta(t, AP, 0.07759873, 0.0000001)

}
