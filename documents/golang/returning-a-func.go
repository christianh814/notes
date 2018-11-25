package main

import "fmt"

func bar() func() int {
	return func() int{
		return 451
	}
}

func main() {
	// below you're assigning f to the function bar. you're essentially loading the function into a variable here
	f := bar()
	// Printing out the type...you see it's a type function that returns an int
	fmt.Printf("%T\n", f)
	// You need to run the actual function to get the return...you can load this to variable
	i := f()
	// This actually prints out 42
	fmt.Println(i)
	// The below here is the same as you ran above
	//	fmt.Println(f())
}
