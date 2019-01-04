package main

import (
	"fmt"
)

type person struct {
	Fname string
	Lname string
}

func (p *person) speak() {
	fmt.Println("Hello")
}

type human interface {
	speak()
}

func saySomething(h human) {
	h.speak()
}

func main() {
	p1 := person{
		Fname: "Christian",
		Lname: "Hernandez",
	}
	//This won't work
	/// saySomething(p1)

	// This will work
	saySomething(&p1)
}

//
//-30-
