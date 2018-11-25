package main

import "fmt"

func main() {
	//break when x is 10
	x := 0
	for {
		fmt.Println(x)
		if x == 10 {
			break
		} else {
			x++
		}
	}
	fmt.Println("done")
	// continue, "skips" the current itteration and goes to the next
	a := 0
	for {
		a++
		if a > 10 {
			break
		}
		//
		if a%2 != 0 {
			continue
		}
		fmt.Println(a)
	}
}
