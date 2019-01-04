package main

import (
	"fmt"
)

func main() {
	// A channel that I can put integers. putting the "1" means it's a buffered channel. You
	// are saying this channel can hold one value
	c := make(chan int, 1)

	// I'm putting "42" into channel "c".
	c <- 42

	// EXAMPLE: the below won't run if you uncomment it because
	// above you specified "1"...meaning only one value can be loaded
	/**
	c <- 57
	**/

	// Load what's in channel "c" into the print statement
	fmt.Println(<-c)
}
