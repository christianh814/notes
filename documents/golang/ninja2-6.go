package main

import "fmt"

const (
	x = iota + 2015
	y
	z
	a
)

func main() {
	fmt.Println(x, y, z, a)
}
