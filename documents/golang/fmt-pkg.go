package main

import "fmt"

var y int = 42

func main() {
	fmt.Println(y)
	fmt.Printf("%T\n", y)
	fmt.Printf("%b\n", y)
	fmt.Printf("%x\n", y)
	fmt.Printf("%#x\n", y)

	y = 911

	// \t means "tab"
	// fmt.Printf("%b\t%x\t%#x\t", y, y, y)
	//fmt.Printf("%#x\a", y)

	// Sprint (string print) lets you assign these to a var
	s := fmt.Sprintf("%T\n", y)
	fmt.Println(s)
}
