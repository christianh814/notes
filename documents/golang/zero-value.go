package main

import "fmt"

var y string
var z int
var w float64

func main() {
	fmt.Println(`The value of "y" is`, y)
	fmt.Printf("%T\n", y)
	y = "Bond, James Bond"
	fmt.Println(`The value of "y" is`, y)

	fmt.Println(`The value of "z" is`, z)
	fmt.Printf("%T\n", z)

	fmt.Println(`The value of "w" is`, w)
	fmt.Printf("%T\n", w)
}
