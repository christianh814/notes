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
}
