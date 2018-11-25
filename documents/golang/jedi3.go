package main

import "fmt"

var x int
var y string
var z bool

func main() {
	x = 42
	y = "James Bond"
	z = true

	s := fmt.Sprintf("The number is %v the name is %v and the value is %v", x, y, z)

	fmt.Println(s)
}
