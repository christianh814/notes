package main

import "fmt"

// A `struct` is a data structure of different types
type person struct {
	first string
	last  string
	age   int
}

//embed a struct inside a struct :-0
type secretAgent struct {
	// this struct has the same struct as a person
	person
	// PLUS a licence to kill
	ltk bool
}

func main() {
	sa1 := secretAgent{
		person: person{
			first: "James",
			last:  "Bond",
			age:   32,
		},
		ltk: true,
	}

	fmt.Println(sa1.first, sa1.last, sa1.age, sa1.ltk)

}
