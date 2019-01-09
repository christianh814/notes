package main

import (
	"fmt"
)

func main() {
	c := make(chan int)    // both sends/recv
	cr := make(<-chan int) // recv (pull off from channel)
	cs := make(chan<- int) //send (load into channel)

	fmt.Printf("c\t%T\n", c)
	fmt.Printf("cr\t%T\n", cr)
	fmt.Printf("cs\t%T\n", cs)
}
