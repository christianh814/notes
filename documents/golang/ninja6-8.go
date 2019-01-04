package main

import "fmt"

func foo() func() {
	return func() {
		fmt.Println("We got the func")
	}
}

func main() {
	x := foo()
	x()
}
