package main

import "fmt"

func foo() int {
	return 2
}

func bar() (int, string) {
	return 5, "Hello"
}

func main() {
	x := foo()
	y, z := bar()
	//
	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)
}
