package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")

	// test the catch of panic
	res, err := TestCatchPanic(10.0, "abc")
	if err != nil {
		fmt.Println("Panic caught:", err)
	}
	fmt.Println("Result:", res)
	

}

// func to test the catch of panic
func TestCatchPanic(a ,b interface{})(res float64,err error){
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	res = a.(float64) / b.(float64)
	return res, err
}