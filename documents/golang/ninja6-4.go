package main

import "fmt"

type person struct {
	fn string
	ln string
	ag int
}

func (p person) speak() {
	fmt.Println("Hello, I am", p.fn, p.ln, ", My age is ", p.ag)
}

func main() {
	p1 := person {
		fn: "Christian",
		ln: "Hernandez",
		ag: 36,
	}
	p1.speak()
}
