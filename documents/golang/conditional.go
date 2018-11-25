package main

import "fmt"

func main() {
	if true {
		fmt.Println("Hello")
	}
	if false {
		fmt.Println("Hello no")
	}
	if !true {
		fmt.Println("Hello no")
	}
	if !false {
		fmt.Println("Hello")
	}
	// initialization statement
	if x := 42; x == 2 {
		fmt.Println("This is false")
	}
	// With else and such
	if y := 47; y == 40 {
		fmt.Println("Value is 40")
	} else if y == 47 {
		fmt.Println("Value is 47")
	} else {
		fmt.Println("Value is unknown")
	}
}
