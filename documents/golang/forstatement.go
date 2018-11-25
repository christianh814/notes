package main

import "fmt"

func main() {
	// Below is effectively a "while" statement
	a := 2
	b := 1024
	for a < b {
		a *= 2
		fmt.Println(a)
	}
}
