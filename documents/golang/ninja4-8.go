package main

import "fmt"

func main() {
	// Here we are mapping a key of a string with a value of a slice of string (that's why the values start with `[]string`
	// Instead of `map[string]string{}` (which stores key, value as an array)...this stores like keys with values of a slice (i.e. arrays)
	m := map[string][]string{
		"bond_james": []string{"Marinis", "women"},
		"penny_miss": []string{"scotch", "men"},
		"hernandez_liz": []string{"takis", "starbux"},
	}
	for k, v := range m {
		fmt.Println("Record for", k)
			for i, v2 := range v {
				fmt.Println("\t", i, v2)
			}
	}
}
