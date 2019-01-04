package main

import "fmt"

func main() {
	//assign a function to a variable
	f1 := func() {
		fmt.Println("Hello func expression")
	}
	f1()
	//
	f2 := func(x int) {
		fmt.Println(x)
	}
	f2(47)
}
