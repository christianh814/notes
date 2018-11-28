package main

import "fmt"
//Closure is limiting the scope of a variable. Below the scope of x is this WHOLE app. Meaning that even functions will get this
var x int

func main() {
	x = 1
	fmt.Println(x)
	foo()
	fmt.Println(x)
	// This  var y is scopped within main() so foo() doesn't know about it
	y := 5
	fmt.Println(y)
	// Starting a new codeblock scopes it within the block
	{
		z := 42
		fmt.Println(z)
	}
	// this won't work...it's out of scope
	//fmt.Println(z)

	fmt.Println("-----------------------------")

	// Let's call increment()
	a := increment()
	b := increment()
	fmt.Println(a())
	fmt.Println(a())
	fmt.Println(a())
	fmt.Println(a())
	fmt.Println(b())
	fmt.Println(b())

}

func foo() {
	x++
	// this won't work
	//y++
}

// function increment() returns a func() that returns an int
func increment() func() int {
	// the var x is scopped within this function
	var x int
	return func() int {
		x++
		return x
	}
}
