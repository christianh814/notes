package main

import (
	"fmt"
	"sort"
)

func main() {
	xi := []int{99, 75, 103, 4, 22, 47, 100000}
	xs := []string{"bob", "mike", "winston", "abby", "shelly", "erica"}

	//unsorted
	fmt.Println(xi)
	fmt.Println(xs)

	//formatting stuff
	fmt.Println("---------------")

	// sorted `sort.Ints` sorts integers
	sort.Ints(xi)
	fmt.Println(xi)

	//sorted `sort.Strings` storts strings
	sort.Strings(xs)
	fmt.Println(xs)
}
