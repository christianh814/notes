package main

import "fmt"

// this function foo() takes a pointer to an int and stores it as y
func foo(y *int) {
	fmt.Println(y)
	fmt.Println(*y)
	*y = 43
	fmt.Println(y)
	fmt.Println(*y)
}

func main() {
	x := 0
	// &x means "pointer to the value of x"
	fmt.Println("x Before", &x)
	fmt.Println("x Before", x)
	foo(&x)
	fmt.Println("x after", &x)
	fmt.Println("x after", x)
}
