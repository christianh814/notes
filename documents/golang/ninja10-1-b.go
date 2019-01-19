package main

import (
	"fmt"
)

func main() {
	// another solution is to make a 'buffered channel'; this "blocks" the memory while it's
	// being written
	c := make(chan int, 1)

	c <- 42

	fmt.Println(<-c)
}
