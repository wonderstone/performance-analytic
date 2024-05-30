package statistics

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

// define ols struct to store the Ordinary Least Squares model
type OLS struct {
	// X is the matrix of independent variables
	X *mat.Dense
	// Y is the matrix of dependent variables
	Y *mat.Dense

	// coefficients is the matrix of coefficients
	coefficients *mat.Dense
	// standard errors is the matrix of standard errors
	standardErrors *mat.Dense
	// tStats is the matrix of t statistics
	tStats *mat.Dense
	// pValues is the matrix of p values
	pValues *mat.Dense
	// residuals is the matrix of residuals
	residuals *mat.Dense
	// yhat is the matrix of fitted values
	yhat *mat.Dense
	// rSquared is the coefficient of determination
	rSquared float64
	// adjRSquared is the adjusted coefficient of determination
	adjRSquared float64

}

// NewOLS creates a new Ordinary Least Squares model
func NewOLS(X, Y *mat.Dense) *OLS {
	return &OLS{
		X: X,
		Y: Y,
	}
}

// ReadCSV reads a csv file and returns the X and Y matrices
func ReadCSV(filename string, Intercept bool, y string, Xs ...string) (*mat.Dense, *mat.Dense) {
	// read the csv file
	csvFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	// create a new csv reader
	reader := csv.NewReader(csvFile)
	// read the header
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}
	// fmt.Println(header)

	// read all the records
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	// fmt.Println(records)

	// check if the independent variable is in the header
	// and store the index in a tmp slice
	var yIndex int
	var xIndex []int
	for x := range Xs {
		found := false
		for i, h := range header {
			if h == Xs[x] {
				found = true
				xIndex = append(xIndex, i)
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("column %s not found in the csv file", Xs[x]))
		}
	}
	for i, h := range header {
		if h == y {
			yIndex = i
			break
		}
	}

	// create the X and Y matrices
	X := mat.NewDense(len(records), len(Xs)+1, nil)
	if Intercept {
		for i := 0; i < len(records); i++ {
			X.Set(i, 0, 1)
			for j, x := range xIndex {
				val, err := strconv.ParseFloat(records[i][x], 64)
				if err != nil {
					panic(err)
				}
				X.Set(i, j+1, val)
			}
		}
	} else {
		for i := 0; i < len(records); i++ {
			for j, x := range xIndex {
				val, err := strconv.ParseFloat(records[i][x], 64)
				if err != nil {
					panic(err)
				}
				X.Set(i, j, val)
			}
		}
	}
	Y := mat.NewDense(len(records), 1, nil)
	for i := 0; i < len(records); i++ {
		val, err := strconv.ParseFloat(records[i][yIndex], 64)
		if err != nil {
			panic(err)
		}
		Y.Set(i, 0, val)
	}
	return X, Y
}

// method Run will do all the OLS model calculations
func (ols *OLS) Run() {
	// calculate the matrix Transpose
	xT := mat.DenseCopyOf(ols.X.T())
	// calculate the matrix xT * x
	xTx := mat.NewDense(ols.X.RawMatrix().Cols, ols.X.RawMatrix().Cols, nil)
	xTx.Product(xT, ols.X)

	// calculate the inverse of xTx
	xTxInv := mat.NewDense(ols.X.RawMatrix().Cols, ols.X.RawMatrix().Cols, nil)
	xTxInv.Inverse(xTx)

	// calculate the matrix xT * y
	xTy := mat.NewDense(ols.X.RawMatrix().Cols, 1, nil)
	xTy.Product(xT, ols.Y)

	// calculate the matrix beta
	ols.coefficients = mat.NewDense(ols.X.RawMatrix().Cols, 1, nil)
	ols.coefficients.Product(xTxInv, xTy)

	// calculate the matrix yHat
	ols.yhat = mat.NewDense(ols.Y.RawMatrix().Rows, 1, nil)
	ols.yhat.Product(ols.X, ols.coefficients)

	// calculate the matrix residuals
	ols.residuals = mat.NewDense(ols.Y.RawMatrix().Rows, 1, nil)
	ols.residuals.Sub(ols.Y, ols.yhat)

	// calculate the σ^2 which is the variance of the residuals
	// turn residuals into a mat.Vector
	residualsVec := mat.NewVecDense(ols.Y.RawMatrix().Rows, nil)
	for i := 0; i < ols.Y.RawMatrix().Rows; i++ {
		residualsVec.SetVec(i, ols.residuals.At(i, 0))
	}
	// calculate the σ^2 which is the variance of the residualsVec
	sigma2 := mat.Dot(residualsVec, residualsVec) / float64(ols.Y.RawMatrix().Rows-ols.X.RawMatrix().Cols)

	// calculate the matrix of standard errors
	ols.standardErrors = mat.NewDense(ols.X.RawMatrix().Cols, 1, nil)
	for i := 0; i < ols.X.RawMatrix().Cols; i++ {
		ols.standardErrors.Set(i, 0,  math.Sqrt(sigma2*xTxInv.At(i, i)))
	}

	// calculate the matrix of t statistics
	ols.tStats = mat.NewDense(ols.X.RawMatrix().Cols, 1, nil)
	for i := 0; i < ols.X.RawMatrix().Cols; i++ {
		ols.tStats.Set(i, 0, ols.coefficients.At(i, 0)/ols.standardErrors.At(i, 0))
	}

	// calculate the matrix of p values
	// calculate the degrees of freedom
	df := float64(ols.Y.RawMatrix().Rows - ols.X.RawMatrix().Cols)
	ols.pValues = mat.NewDense(ols.X.RawMatrix().Cols, 1, nil)
	for i := 0; i < ols.X.RawMatrix().Cols; i++ {
		dist := distuv.StudentsT{Mu: 0, Sigma: 1, Nu: df}
		ols.pValues.Set(i, 0, 2*(1-dist.CDF(math.Abs(ols.tStats.At(i, 0)))))
	}
	// calculate the rSquared
	// calculate the sum of squares total
	yMean := mat.Sum(ols.Y) / float64(ols.Y.RawMatrix().Rows)
	ssTotal := mat.NewDense(ols.Y.RawMatrix().Rows, 1, nil)
	yMeanMat := mat.NewDense(ols.Y.RawMatrix().Rows, 1, nil)
	for i := 0; i < ols.Y.RawMatrix().Rows; i++ {
		yMeanMat.Set(i, 0, yMean)
	}
	ssTotal.Sub(ols.Y, yMeanMat)
	ssTotal.MulElem(ssTotal, ssTotal)
	ssTotalVal := mat.Sum(ssTotal)

	// calculate the sum of squares residuals
	ssResiduals := mat.NewDense(ols.residuals.RawMatrix().Rows, 1, nil)
	ssResiduals.MulElem(ols.residuals, ols.residuals)
	ssResidualsVal := mat.Sum(ssResiduals)
	// calculate the rSquared
	ols.rSquared = 1 - ssResidualsVal/ssTotalVal
	// calculate the adjusted rSquared
	ols.adjRSquared = 1 - (1-ols.rSquared)*(float64(ols.Y.RawMatrix().Rows-1)/float64(ols.Y.RawMatrix().Rows-ols.X.RawMatrix().Cols))
}

// method to get the coefficients
func (ols *OLS) Coefficients() *mat.Dense {
	return ols.coefficients
}

// method to get the standard errors
func (ols *OLS) StandardErrors() *mat.Dense {
	return ols.standardErrors
}

// method to get the t statistics
func (ols *OLS) TStats() *mat.Dense {
	return ols.tStats
}

// method to get the p values
func (ols *OLS) PValues() *mat.Dense {
	return ols.pValues
}

// method to get the residuals
func (ols *OLS) Residuals() *mat.Dense {
	return ols.residuals
}

// method to get the fitted values
func (ols *OLS) YHat() *mat.Dense {
	return ols.yhat
}

// method to get the rSquared
func (ols *OLS) RSquared() float64 {
	return ols.rSquared
}

// method to get the adjusted rSquared
func (ols *OLS) AdjRSquared() float64 {
	return ols.adjRSquared
}

// method to get the number of observations
func (ols *OLS) N() int {
	return ols.Y.RawMatrix().Rows
}

// method to get the number of independent variables
func (ols *OLS) K() int {
	return ols.X.RawMatrix().Cols
}

// method to get the degrees of freedom
func (ols *OLS) DF() float64 {
	return float64(ols.Y.RawMatrix().Rows - ols.X.RawMatrix().Cols)
}

