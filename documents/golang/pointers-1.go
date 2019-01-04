package main

import "fmt"

func main() {
	a := 42
	// the & before the variable shows it's address in memory
	fmt.Println(&a)

	fmt.Printf("%T\n", a)
	// this will return *int which means "it's a pointer to an int"
	fmt.Printf("%T\n", &a)

	// the * gives you the value that's stored at an address
	fmt.Println(*&a)
}
