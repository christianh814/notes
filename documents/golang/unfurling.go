package main

import "fmt"

// store an unlimited amount of integers in x
func foo(x ...int) int {
	sum := 0
	for _, v := range x {
		sum += v
	}
	return sum
}

func main() {
	xi := []int{1, 1, 1, 1, 1}
	// like FM, you "unfurl" the slice into the function
	y := foo(xi...)
	fmt.Println(y)
}
