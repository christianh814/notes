package main

import "fmt"

// A `struct` is a data structure of different types
type person struct{
	first string
	last string
	age int
}

func main() {
	p1 := person {
		first: "Christian",
		last: "Hernandez",
		age: 36,
	}
	p2 := person {
		first: "Erica",
		last: "Macha",
		age: 33,
	}

	fmt.Println(p1, p2)

	// access a specific field
	fmt.Println(p1.first)
	fmt.Println(p2.age)
}
