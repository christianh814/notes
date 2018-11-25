package main

import "fmt"

func main() {
	x := 42
	fmt.Printf("%d\n", x)
	fmt.Printf("%b\n", x)
	fmt.Printf("%#x\n", x)
	// OR
	fmt.Printf("%d\t%b\t%#x\n", x, x, x)
}
