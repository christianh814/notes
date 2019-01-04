package main

import "fmt"

// use loops instead of recursion...this is just here for your info

// the function factorial() takes an argument as an int and stores it in the var n; factorial() also returns an int
func factorial(n int) int {
	// you need some way to stop
	if n == 0 {
		return 1
	}
	// returns the number then calls itself again to multiply that number -1 adnausum
	return n * factorial((n - 1))
}

func main() {
	num := factorial(4)
	yum := loopFac(4)
	fmt.Println(num)
	fmt.Println(yum)
}

/*
This is how you would do it with a loop
*/
func loopFac(n int) int {
	total := 1
	for ; n > 0; n-- {
		total *= n
	}
	return total
}
