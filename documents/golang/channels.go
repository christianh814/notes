package main

import (
	"fmt"
)

func main() {
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
