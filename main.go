package main

import (
	"fmt"
	"math"
	"wonderstone/performance-analytics/statistics"
)

// define Evolutionary direction struct

type EvolveDirect interface {
	// take []float64 as input and return float64
	// but make sure fit the GEP regime 
	// say "positive and the more the better"
	EvoDct([]float64) float64
}



// give a struct to implement the interface
type ED struct {
	// in case may need some parameters
}

// implement the interface
func (ed ED) EvoDct(Ra []float64) float64 {
	tmpRes := statistics.SharpeRatio(Ra, 0.035/12, 12, true)
	return sigmoid(tmpRes)
}


// say sigmoid function
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}


// give a function to use the interface
func useInterface(ed EvolveDirect, ra []float64) {
	result := ed.EvoDct(ra)
	fmt.Printf("The result is %f\n", result)
}


func main() {
	fmt.Println("Hello, playground")
	// test the catch of error
	res, err := TestCatchError(10.0, "abc")
	if err != nil {
		fmt.Println("Error caught:", err)
	}
	fmt.Println("Result:", res)
	// test the interface
	// get the returns
	// define the returns
	dt, fds := statistics.ReadData("./data/managers.csv")


	rtp, _ := statistics.CheckPos(fds, "HAM1")
	rts := statistics.GetSecondDimensionData(dt, rtp)
	rt, e := statistics.TryStringToFloatSlice(rts)
	if e != nil {
		panic(e)
	}

	// use a struct to implement the interface
	tmpED := ED{}
	useInterface(tmpED, rt)


}

// func to test the catch of Error
func TestCatchError(a ,b interface{})(res float64,err error){
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	res = a.(float64) / b.(float64)
	return res, err
}