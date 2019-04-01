package main

import (
	"fmt"
	"./greaterthan"
)

func main() {
	a := 10
	b := 2
	if GreaterThan(a, b) {
		fmt.Printf("Number %v is greater than %v\n", a, b)
	}
}
