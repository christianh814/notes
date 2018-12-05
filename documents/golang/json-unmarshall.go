package main

import (
	"fmt"
	"encoding/json"
)

type person struct {
	// Here you're saying take the json field and convert it to the key "First" in this struct. When I recv a JSON array I need a place to store it...this is that place
	First string `json:"First"`
	Last string `json:"Last"`
	Age int `json:"Age"`
}

func main() {
	// This function takes a slice of byte and an interface and returns an error
	// func Unmarshall(data []byte, v interface{}) (err)
	s := `[{"First":"Christian","Last":"Hernandez","Age":36},{"First":"Wendy","Last":"Thomas","Age":39}]`
	bs := []byte(s)

	// This decalres and emtpy var of people of type person...which is the struct we are using to store the json
	//people := []person{}
	var people []person


	// this takes the byte slice and the address pointer to people. It basically loads the json into the struct that is people
	err := json.Unmarshal(bs, &people)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(people)

	for i, v := range people {
		fmt.Println("\tPERSON NUMBER", i)
		fmt.Println(v.First)
		fmt.Println(v.Last)
		fmt.Println(v.Age)
	}
}
