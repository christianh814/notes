package main

import "fmt"

func main() {
	x := 42
	if x == 42 && x == 41 {
		fmt.Println("this shouldn't print")
	}
	if x == 41 || x == 42 {
		fmt.Println("this WILL print")
	}
}
