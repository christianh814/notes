package main

import (
	"fmt"
)

func main() {
	// create a channel with two intergers loaded.
	ch := make(chan int, 2)
	// load the two integers
	ch <- 42
	ch <- 43
	fmt.Println(<-ch) //load the first entry in the channel into this print statement
	fmt.Println(<-ch) //load the second entry in the channel into this print statement
	// Channels allow us to pass values between goroutines


	// A channel that I can put integers. 
	c := make(chan int)

	go func() {
		// I'm putting "42" into channel "c".
		// Loading into channels must be done in it's own goroutine
		c <- 42
	}()

	// Load what's in channel "c" into the print statement
	fmt.Println(<-c)
}
