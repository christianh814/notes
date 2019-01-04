package main

import "fmt"

func foo() int {
	var x int
	x = 42
	return x
}

func main() {
	g := foo()
	fmt.Println(g)
}
