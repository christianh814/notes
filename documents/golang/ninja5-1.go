package main

import "fmt"

type person struct {
	first   string
	last    string
	flavors []string
}

func main() {
	p1 := person{
		first: "Christian",
		last:  "Hernandez",
		flavors: []string{
			"chocolate",
			"vanilla",
			"strawberry",
		},
	}
	fmt.Println(p1.first)
	//
	for i, v := range p1.flavors {
		fmt.Println("The index is:", i, "And the value is:", v)
	}
}
