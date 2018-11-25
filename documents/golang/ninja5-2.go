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

	p2 := person{
		first: "Micheal",
		last:  "Hodges",
		flavors: []string{
			"everything but the cow",
			"brownie core",
			"half baked",
		},
	}

	fmt.Println(p1.first)
	//
	for i, v := range p1.flavors {
		fmt.Println("The index is:", i, "And the value is:", v)
	}
	// you can make a map of your struct! holy cow!
	m := map[string]person{
		p1.last: p1,
		p2.last: p2,
	}
	fmt.Println(m["Hernandez"])
	//
	for k, v := range m {
		fmt.Println(k, "=>", v.first, v.last)
		for _, val := range v.flavors {
			fmt.Printf("\tA Favorate flavor: %v\n", val)
		}
	}
}
//
//
