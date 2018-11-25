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
	y := foo(1, 1, 1, 1, 1)
	fmt.Println(y)
}
