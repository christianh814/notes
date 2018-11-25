package main

import "fmt"

func main() {
	x := map[string]int{
		"Dog": 5,
		"Cat": 55,
		"Mouse": 88,
		"Bird": 99,
	}
	fmt.Println(x)
	// Delete stuff from map syntax is `delete(mapname, "key")`
	delete(x, "Bird")
	fmt.Println(x)
	// practicing for because why not
	for i := 0; i <= 5; i++ {
		fmt.Println(i)
	}
}
