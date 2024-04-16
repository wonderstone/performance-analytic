package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test CAPM all methods
func TestCAPMAllMethods(t *testing.T) {
	dt, fds := ReadData("../data/managers.csv")
	// define the returns for HAM1
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
	// define the rf return
	rf, _ := CheckPos(fds, "US 3m TR")
	rfs := GetSecondDimensionData(dt, rf)
	r, e := TryStringToFloatSlice(rfs)
	if e != nil {
		panic(e)
	}
	capm := CAPM{rt, bm}

	// test the Beta method
	beta := capm.Beta()
	// only check the first 7 decimal places
	assert.InDelta(t, 0.3906033, beta, 0.0000001, "Beta should be 0.3906033")

	// test the alpha method
	alpha := capm.Alpha(0.0)
	// only check the first 9 decimal places
	assert.InDelta(t, 0.007738016, alpha, 0.000000001, "Alpha not equal to 0.007738016")
	alpha = capm.Alpha(0.04/12)
	// only check the first 9 decimal places
	assert.InDelta(t, 0.005706694, alpha, 0.000000001, "Alpha not equal to 0.005706694")
	// ! again this part is different from the original PA package
	// ! choose on your own risk
	// ! the original value is 0.005774729
	// ! the value here is 0.005771834859331053
	alpha = capm.Alpha(r)
	assert.InDelta(t, 0.005771835, alpha, 0.000000001, "Alpha not equal to 0.005771835")
	// test the TimingRatio method
	tr := capm.TimingRatio()
	// only check the first 7 decimal places
	assert.InDelta(t, 0.7070631, tr, 0.0000001, "TimingRatio should be 0.7070631")

}
