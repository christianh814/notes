package main

import "fmt"

func foo(x ...int) int {
	total := 0
	for _, v := range x {
		total += v
	}
	return total
}

func bar(xi []int) int {
	sum := 0
	for _, v := range xi {
		sum += v
	}
	return sum
}

func main() {
	x := []int{1, 2, 3, 4}
	fmt.Println(foo(x...))
	//
	y := []int{1, 2, 3, 4}
	fmt.Println(bar(y))
}
