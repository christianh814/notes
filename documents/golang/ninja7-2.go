package main

import "fmt"

type person struct {
	fn string
	ln string
	ag int
}

// this takes an argument of the value of type person....see below...
func changeMe(p *person) {
	// Changing the fn only
	p.fn = "Adrian"
}

func main() {
	p1 := person{
		fn: "Christian",
		ln: "Hernandez",
		ag: 36,
	}
	fmt.Println(p1)
	// We are passing in the mem address which then the func "de-references" it with the *
	changeMe(&p1)
	fmt.Println(p1)
}
