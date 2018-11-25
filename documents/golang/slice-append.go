package main

import "fmt"

func main() {
	x := []string{"cat", "dog", "mouse"}
	fmt.Println(x)
	// use the keyword `append`. basic syntax is `append(sliceVar, element1, element2, ...)`
	x = append(x, "bird")
	fmt.Println(x)

	// make a big array from two smaller ones????
	y := []string{"tiger", "lion", "bear"}
	// here the `...` after the y means "take all the values from y and put them here"
	x  = append(x, y...)
	fmt.Println(x)
}
