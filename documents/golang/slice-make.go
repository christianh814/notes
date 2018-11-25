package main

import "fmt"

func main() {
	// Slice is built on an array. Slice is dynamic whereas arrays are not. So when you manipulate a slice a new array needs to be made in memory.
	// To middigate that; you can use `make` if you know how many elements you're going to store upfront for effeciency.
	// make takes "type" "lenth" and "capacity"
	x := make([]int, 10, 100)
	fmt.Println(x)
	fmt.Println(len(x))
	fmt.Println(cap(x))

	// you can assign values to specific indexes like this
	x[4] = 42
	fmt.Println(x)

	// BUT only up to 10...if you want to do index 10; you need to "append"
	x = append(x, 47)
	fmt.Println(x)
	fmt.Println(len(x))
	fmt.Println(cap(x))

	// You can add beyond your cap; but what happens is that the original underlying array is thrown away; and a new one is made to accomidate.
	// This is not good practice but possible
}
