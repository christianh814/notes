package main

import "fmt"

func sum(xi ...int) int {
	total := 0
	for _, v := range xi {
		total += v
	}
	return total
}

// CALL BACK IS A function THAT TAKES ANOTHER function AS AN ARGUMENT

// this is an example of a callback. the `f` is expecting to be passed in `func(xi ..int) int` and `ix` is expecting 0 to infiniate integers (slice of int)
func even(f func(xi ...int) int, ix ...int) int {
	// decalring `yi` as a slice of int
	var yi []int
	// range over the slice that was passed with `ix` and add only the even numbers to  the empty slice of int that is `yi`
	for _, v := range ix {
		if v%2 == 0 {
			yi = append(yi, v)
		}
	}
	//return the values and returning an `int`...which is what you said this fucntion will return
	return f(yi...)
}

// ODD is similar
func odd(f func(xi ...int) int, ix ...int) int {
	// decalring `yi` as a slice of int
	var yi []int
	// range over the slice that was passed with `ix` and add only the even numbers to  the empty slice of int that is `yi`
	for _, v := range ix {
		if v%2 != 0 {
			yi = append(yi, v)
		}
	}
	//return the values and returning an `int`...which is what you said this fucntion will return
	return f(yi...)
}

func main() {
	ii := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := sum(ii...)
	fmt.Println("All numbers summed up", s)
	//
	j := even(sum, ii...)
	fmt.Println("All EVEN Numbers summed up", j)
	//
	o := odd(sum, ii...)
	fmt.Println("All ODD Numbers summed up", o)

}
