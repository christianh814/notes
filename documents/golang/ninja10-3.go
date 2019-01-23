package main

import (
	"fmt"
)

func main() {
	c := gen()
	receive(c)

	fmt.Println("about to exit")
}

func receive(c <-chan int) {
	for i := range c {
		fmt.Println(i)
	}
}

func gen() <-chan int {
	c := make(chan int)

	go func() {
		for i := 0; i < 11; i++ {
			c <- i
		}
		close(c)
	}()

	return c
}
