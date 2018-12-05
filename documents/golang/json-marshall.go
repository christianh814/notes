package main

import (
	"fmt"
	"encoding/json"
)

// In order for json.Marshal to work; the fields need to be upper case
type person struct {
	First string
	Last string
	Age int
}

func main() {
	p1 := person {
		First: "Christian",
		Last: "Hernandez",
		Age: 36,
	}
	p2 := person {
		First: "Wendy",
		Last: "Thomas",
		Age: 39,
	}

	// I'm creating an array. "people" is a "slice" of "person" with p1 and p2
	people := []person{p1, p2,}

	// This function takes an interface and returns a slice of byte and an error
	// func Marshal(v interface{}) ([]byte, err)
	bs, err := json.Marshal(people)
	// if error is NOT empty; it means we got an error...so print it
	if err != nil {
		fmt.Println(err)
	}
	// convert your byte slice to a string and print it
	fmt.Println(string(bs))
}
