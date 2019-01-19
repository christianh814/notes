package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

	// first solution is to just launch the channel
	// into it's own goroutine with an anonymous function
	go func() {
		c <- 42
	}()

	fmt.Println(<-c)
}
