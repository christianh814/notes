package main

import (
	"fmt"
)

func main() {
	// create a channel that can ONLY send
	c := make(chan<- int, 2)

	c <- 42
	c <- 43
	// These will fail
	fmt.Println(<-c)
	fmt.Println(<-c)

	fmt.Println("------")
	fmt.Printf("%T\n", c)
}
