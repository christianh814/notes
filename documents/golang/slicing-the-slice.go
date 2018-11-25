package main

import "fmt"

func main() {
	x := []int{0, 1, 2, 3, 4, 5}
	fmt.Println(x)
	// A : can be used to get a range
	//
	// This prints from position 2 til the end
	fmt.Println(x[1:])
	// This prints from pos 2 through 4
	fmt.Println(x[2:4])
}
