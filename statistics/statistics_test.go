package statistics

// test the ActivePremium function using
// testify package
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// * function to change []string to []float64
var dt [][]string
var fds []string

// ~ init this test
func init() {
	dt, fds = ReadData("../data/managers.csv")
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
	rt1, bm1 := StringToFloatSliceBench(rts1, bms)

	// calculate the active premium
	// this number（0.07759873）is from the R code
	AP = ActivePremium(rt1, bm1, 12, true)
	assert.InDelta(t, AP, 0.07759873, 0.0000001)

}
