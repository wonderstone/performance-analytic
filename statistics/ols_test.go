package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test ols ReadCSV
func TestReadCSV(t *testing.T) {
	// read the csv file
	X, Y := ReadCSV("data.csv", true, "index_price", "interest_rate", "unemployment_rate")
	// check if the X matrix is correct
	assert.Equal(t, 24, X.RawMatrix().Rows)
	assert.Equal(t, 3, X.RawMatrix().Cols)
	// check if the Y matrix is correct
	assert.Equal(t, 24, Y.RawMatrix().Rows)
	assert.Equal(t, 1, Y.RawMatrix().Cols)

	// Run the OLS model
	ols := NewOLS(X, Y)
	ols.Run()

	// check if the Coefficients matrix is correct
	Coefficients := ols.Coefficients()
	// check the three values are （1798.4039776248844，345.5400870107158，-250.14657136921616）
	assert.InDelta(t, 1798.4039776248844, Coefficients.At(0, 0), 0.0001)
	assert.InDelta(t, 345.5400870107158, Coefficients.At(1, 0), 0.0001)
	assert.InDelta(t, -250.14657136921616, Coefficients.At(2, 0), 0.0001)

	// check if the  standard error is correct
	standardErrors := ols.StandardErrors()
	// check the first three values are ⎡  899.248074996123⎤ ⎢111.36692222534533⎥ ⎣117.94986891695314⎦
	assert.InDelta(t, 899.248074996123, standardErrors.At(0, 0), 0.0001)
	assert.InDelta(t, 111.36692222534533, standardErrors.At(1, 0), 0.0001)
	assert.InDelta(t, 117.94986891695314, standardErrors.At(2, 0), 0.0001)	

	// check if the tStats matrix is correct
	tStats := ols.TStats()
	// check the first three values are ⎡1.9998975006231043⎤ ⎢ 3.102717396746701⎥ ⎣-2.120787192610963⎦
	assert.InDelta(t, 1.9998975006231043, tStats.At(0, 0), 0.0001)
	assert.InDelta(t, 3.102717396746701, tStats.At(1, 0), 0.0001)
	assert.InDelta(t, -2.120787192610963, tStats.At(2, 0), 0.0001)

	// check if the pValues matrix is correct
	pValues := ols.PValues()
	// check the first three values are ⎡  0.05861188624935032⎤ ⎢0.0053891736105109445⎥ ⎣  0.04601273885306845⎦
	assert.InDelta(t, 0.05861188624935032, pValues.At(0, 0), 0.0001)
	assert.InDelta(t, 0.0053891736105109445, pValues.At(1, 0), 0.0001)
	assert.InDelta(t, 0.04601273885306845, pValues.At(2, 0), 0.0001)

	// check if the rsquared is correct
	rSquared := ols.RSquared()
	assert.InDelta(t, 0.8976335894170216,rSquared, 0.0001)
	// check if the adjRSquared is correct
	adjRSquared := ols.AdjRSquared()
	assert.InDelta(t, 0.887884407456738, adjRSquared, 0.0001)
}