package main

import (
	"fmt"
)

func main() {
	// Changed it to a bidirectional channel
	//cs := make(chan<- int)
	cs := make(chan int)

	go func() {
		cs <- 42
	}()
	fmt.Println(<-cs)

	fmt.Printf("------\n")
	fmt.Printf("cs\t%T\n", cs)
}

