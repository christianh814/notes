package main

import "fmt"

func main() {
	// no expression defaults to true...the bottow is like typing `switch true { ... }`
	switch {
	case 1 == 2:
		fmt.Println("This will not print")
	case 1 == 1:
		fmt.Println("Print")
	}
}
