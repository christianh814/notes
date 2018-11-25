package main

import "fmt"

var a int

type hotdog int

var b hotdog

func main() {
	a = 42
	b = 99
	fmt.Println(a)
	fmt.Println(b)

	// Let's convert it

	a = int(b)
	fmt.Println(a)
}
