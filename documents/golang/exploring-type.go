package main

import "fmt"

var y = 42
var z = `James Said "Shaken not stirred"`

func main() {
	fmt.Printf("%T\n", y)
	fmt.Println(z)
}
