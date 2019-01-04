package main

import "fmt"

//func (r receiver) identifier (parameters) (returns(s)) { ... }
func foo() {
	fmt.Println("foo")
}

// This takes a parameter as TYPE string and assings it to `s` so you can call it from within the function
// EVERYTHING in Golan is "pass by value"
func bar(s string) {
	fmt.Println(s)
}

func bazz(s string) bool {
	return s == "boo"
}

// take multiples
func name(fn, ln string) map[string]string {
	return map[string]string{
		"First": fn,
		"Last":  ln,
	}
}

//

func main() {
	foo()
	bar("Hello bar")
	fmt.Println(bazz("boo"))
	x := name("Christian", "Hernandez")
	fmt.Println(x["First"])
}
