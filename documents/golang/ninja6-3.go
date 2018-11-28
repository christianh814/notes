package main

import "fmt"

func foo() {
	fmt.Println("Hello")
}

func bar() {
	fmt.Println("World")
}

func main() {
	defer foo()
	bar()
}
