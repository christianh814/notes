package main

import (
	"encoding/json"
	"fmt"
)

//preparing this struct to recv a json array (hits at google with "go to json")
type person struct {
	First string `json:"First"`
	Age   int    `json:"Age"`
}

func main() {
	j := `[{"First":"James","Age":47},{"First":"Ricky","Age":22},{"First":"Peter","Age":77}]`
	var people []person

	//`j` is coming in as a string and Unmarshal expects a slice of byte and the memmory address of people
	err := json.Unmarshal([]byte(j), &people)

	if err != nil {
		fmt.Println(err)
	}

	// people is now a slice
	fmt.Println(people)

	for i, v := range people {
		fmt.Println("Person #", i)
		fmt.Println("\tFirst Name:", v.First)
		fmt.Println("\tAge:", v.Age)
	}
}
