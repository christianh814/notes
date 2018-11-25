package main

import "fmt"

func main() {
	if 1 == 2 {
		fmt.Println("This does not print")
	} else if 1 == 3 {
		fmt.Println("This does not print")
	} else {
		fmt.Println("print")
	}
}
