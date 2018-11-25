package main

import "fmt"

func main() {
	m := map[string][]string{
		"bond_james": []string{"Marinis", "women"},
		"penny_miss": []string{"scotch", "men"},
		"hernandez_liz": []string{"takis", "starbux"},
	}
	// add to it
	delete(m, "bond_james")
	//
	for k, v := range m {
		fmt.Println(k, v)
	}
}
