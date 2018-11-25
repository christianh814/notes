package main

import "fmt"

func main() {
	x := []int{1, 2, 3}
	// context: "i" stands for "index" and "v" stands for "value"
	for i, v := range x {
		fmt.Println(i, v)
	}
	//
	// I am using _ for the index (i.e. usually "i" goes there) because I just care about the values not thier index.
	animals := []string{"cat", "dog", "mouse"}
	for _, pet := range animals {
		fmt.Println(pet)
	}
}
