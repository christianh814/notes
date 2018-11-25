package main

import "fmt"

func main() {
	// when you decalre an array; you have to specify the size (here as [5])...meaning 5 items in this array (0-4)
	var x [5]int
	fmt.Println(x)

	// x at positon 3 is 42 (pos 3 comes 4 because it's a 0 based index)
	x[3] = 42
	fmt.Println(x)
	// prints the length of x
	fmt.Println(len(x))

	// Go specification doc says "don't use arrays...use slices"
}
