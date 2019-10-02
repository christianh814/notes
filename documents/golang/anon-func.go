package main

import "fmt"

func main() {
	func() {
		fmt.Println("Anon func")
	}() // The () here runs it after you've defined it like with `foo()`
	//
	func(x int) {
		fmt.Println("The number is", x)
	}(5) // you pass arguemnts like you would any other function
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}
