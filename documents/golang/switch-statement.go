package main

import "fmt"

func main() {
	// Swich statement are like "case" in shell...kind of
	switch {
	case false:
		fmt.Println("This is false1")
	case true:
		fmt.Println("This is true1")
	case 2 == 2:
		fmt.Println("This is true2")
	case 2 == 3:
		fmt.Println("This is false2")
	}
	// Swich statement; you need to specify fallthrough
	switch {
	case false:
		fmt.Println("This is false1")
	case true:
		fmt.Println("This is true1")
		//in order to print 'true2' you need it to "fallthrough"
		fallthrough
	case 2 == 2:
		fmt.Println("This is true2")
	case 2 == 3:
		fmt.Println("This is false2")
	}
	// default means if nothing is true do this
	switch {
	case false:
		fmt.Println("won't print")
	case false:
		fmt.Println("won't print")
	default:
		fmt.Println("prints")
	}
	// Swich on a value .you can also do `n := "Hi"` and do `switch n { ... }`
	switch "Hi" {
	case "No":
		fmt.Println("Does not print")
	case "Hi":
		fmt.Println("Hell hello there")
	}
	// Multiple values!
	switch "Bruh" {
	case "No":
		fmt.Println("Does not print")
	case "Hi", "Hello", "Bruh":
		fmt.Println("Hell hello there")
	}
}
