package main

import "fmt"

var z int
var w float64
var a int8 = -128
var b int8 = -127

func main() {
	x := 42
	y := 41.999
	z = 57
	w = 5.17
	fmt.Println(x)
	fmt.Println(y)
	fmt.Printf("%T\n", x)
	fmt.Printf("%T\n", y)
	fmt.Println(z)
	fmt.Println(w)
	fmt.Printf("%T\n", z)
	fmt.Printf("%T\n", w)
}
