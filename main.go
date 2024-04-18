package main

import (
	"fmt"
	"math"
	"wonderstone/performance-analytics/statistics"
)

// define a simple Evolutionary direction struct

type EvolveDirect interface {
	// take []float64 as input and return float64
	// but make sure fit the GEP regime
	// say "positive and the more the better"
	EvoDct([]float64) float64
}

// give a struct to implement the interface
// ED??!! Are you kidding me??
type ED struct {
	// * field left in case may need some freaky parameters
	Pars map[string]interface{}
	// * field for result float64? not mandatory
	Result float64
}

// implement the interface
func (ed ED) EvoDct(Ra []float64) float64 {
	tmpRes := statistics.SharpeRatio(Ra, 0.035/12, 12, true)

	// try get target return from Pars
	value, ok := ed.Pars["TargetReturn"]
	if ok {
		return sigmoid(tmpRes,value.(float64))
	}

	return sigmoid(tmpRes,1.0)
}

// for instance: sigmoid function
func sigmoid(x float64, scale float64) float64 {
	return 1 / (1 + math.Exp(-x/scale))
}

// give a function to use the interface
func useInterface(ed EvolveDirect, ra []float64) {
	result := ed.EvoDct(ra)
	fmt.Printf("The result is %f\n", result)
}

func main(){
	fmt.Println("Hello, playground")
	// test the catch of error
	res, err := TestCatchError(10.0, "abc")
	if err != nil {
		fmt.Println("Error caught:", err)
	}
	fmt.Println("Result:", res)
	// test the interface
	// get the returns
	dt, fds := statistics.ReadData("./data/managers.csv")


	// try math.-inf for sigmoid

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
func TestCatchError(a, b interface{}) (res float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	res = a.(float64) / b.(float64)
	return res, err
}
