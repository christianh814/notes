package main

import "fmt"

func foo() {
	fmt.Println("foo")
}

func bar() {
	fmt.Println("bar")
}

func main() {
	// foo() won't run until bar is done running
	defer foo()
	bar()
} //this is where foo() runs...right at the end of the surrounding fucntion...in this case main()
