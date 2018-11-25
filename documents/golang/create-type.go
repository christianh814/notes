package main

import "fmt"

var a int

// You create your own type called "hotdog" and it's underlying type is "int"
type hotdog int

var b hotdog

func main() {
	a = 42
	b = 43
	fmt.Println(a)
	fmt.Println(b)

	// Print type hotdog
	fmt.Printf("%T\n", b)
}
