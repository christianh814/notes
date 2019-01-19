package main

import (
	"fmt"
)

func main() {
	// Made this bidirectional too
	//cr := make(<-chan int)
	cr := make(chan int)

	go func() {
		cr <- 42
	}()
	fmt.Println(<-cr)

	fmt.Printf("------\n")
	fmt.Printf("cr\t%T\n", cr)
}

