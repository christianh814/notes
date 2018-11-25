package main

import "fmt"

func main() {
	m := map[string]int{
		"Paul": 47,
		"Mike": 22,
		"Diane": 59,
	}
	fmt.Println(m)
	// add something
	m["Erica"] = 33
	// range over them
	for k, v := range m {
		fmt.Println("The key is:", k, "and the value is:", v)
	}
}
