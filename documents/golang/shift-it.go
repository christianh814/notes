package main

import "fmt"

func main() {
	x := 2
	fmt.Printf("%d\t\t%b\n", x, x)

	// Shift the binary over one
	y := x << 1
	fmt.Printf("%d\t\t%b\n", y, y)
}
