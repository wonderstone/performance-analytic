package statistics

// test the ActivePremium function using
// testify package
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test ReturnsCalculator all methods
func TestReturnsCalculator(t *testing.T) {

	dt, fds := ReadData("../data/managers.csv")
	// define the returns for HAM1
	rtp, _ := CheckPos(fds, "HAM1")
	rts := GetSecondDimensionData(dt, rtp)
	rt, e := TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}
	rc := ReturnsCalculator{rt}

	// define the benchmark returns
	bmp, _ := CheckPos(fds, "SP500 TR")
	bms := GetSecondDimensionData(dt, bmp)
	bm, e := TryStringToFloatSlice(bms)
	if e != nil {
		panic(e)
	}

	// test the Excess method
	excessReturns:= rc.Excess(0.04 / 12)
	//  only check the first 7 decimal places
	assert.InDelta(t, excessReturns[0], 0.004066667, 0.0000001)
	assert.Equal(t, len(excessReturns), len(rt))

	excessReturns= rc.Excess(bm)
	// only check the first 6 decimal places
	assert.InDelta(t, excessReturns[0], -0.02660, 0.00001)
	assert.Equal(t, len(excessReturns), len(rt))


	// test the Cumulative method
	CumReturn:= rc.Cumulative(true)
	// only check the first 7 decimal places
	assert.InDelta(t, CumReturn, 3.126671, 0.000001)

	// test the StdDev method
	std := StdDev(rt)
	// only check the first 8 decimal places
	assert.InDelta(t, std, 0.02562881, 0.0000001)

	// test the StdDevAnnualized method
	stdAnnualized := StdDevAnnualized(rt, 12)
	// only check the first 7 decimal places
	assert.InDelta(t, stdAnnualized, 0.0887808, 0.000001)


}
