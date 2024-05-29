package statistics

// test the ActivePremium function using
// testify package
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test the common functions
func TestCommon(t *testing.T) {
	// test the skewness function
	dt, fds := ReadData("../data/managers.csv")
	// define the returns for HAM1
	rtp, _ := CheckPos(fds, "HAM1")
	rts := GetSecondDimensionData(dt, rtp)
	rt, e := TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}
	// test the skewness function
	sk := Skewness(rt,"moment")
	// only check the first 7 decimal places
	assert.InDelta(t,  -0.6588445, sk, 0.0000001, "Skewness moment should be -0.6588445")
	sk = Skewness(rt,"sample")
	// only check the first 7 decimal places
	assert.InDelta(t,  -0.6740873, sk, 0.0000001, "Skewness sample should be -0.6740873")
	sk = Skewness(rt,"fisher")
	// only check the first 7 decimal places
	assert.InDelta(t,  0.5695854, sk, 0.0000001, "Skewness fisher should be 0.5695854")



	// test the kurtosis function
	kt := Kurtosis(rt,"moment")
	// only check the first 6 decimal places
	assert.InDelta(t, 5.361589, kt, 0.000001, "Kurtosis moment should be 5.361589")

	kt = Kurtosis(rt,"excess")
	// only check the first 6 decimal places
	assert.InDelta(t, 2.361589, kt, 0.000001, "Kurtosis excess should be 2.361589")

	kt = Kurtosis(rt,"sample")
	// only check the first 6 decimal places
	assert.InDelta(t, 5.570361, kt, 0.000001, "Kurtosis sample should be 5.570361")

	kt = Kurtosis(rt,"sample_excess")
	// only check the first 7 decimal places
	assert.InDelta(t, 2.500415, kt, 0.000001, "Kurtosis should be 2.500415")

	kt = Kurtosis(rt,"fisher")
	// only check the first 7 decimal places
	assert.InDelta(t, 0.8846114, kt, 0.0000001, "Kurtosis fisher should be 0.8846114")
}