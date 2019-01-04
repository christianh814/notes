package main

import "fmt"

func foo() {
	fmt.Println("You got the func")
}

func main() {
	x := foo
	x()
}
