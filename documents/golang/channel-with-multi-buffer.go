package main

import (
	"fmt"
)

func main() {
	// A channel that I can put integers. putting the "2" means it's a buffered channel. You
	// are saying this channel can hold one value
	c := make(chan int, 2)

	// I'm putting "42" into channel "c".
	c <- 42

	// I'm putting "57" into channel "c".
	c <- 57

	// Load what's in the first channel "c" into the print statement
	fmt.Println(<-c)

	// pull the second value sotred in "c" into the print statement
	fmt.Println(<-c)
}
