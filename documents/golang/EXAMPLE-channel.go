package main

import (
	"fmt"
)

func doSomething(x int) int {
	return x * 5
}

func main() {
	// create a channel
	ch := make(chan int)
	go func() {
		//this is a goroutine...this "forks" off but you can save the output in the channel
		ch <- doSomething(5)
	}()
	fmt.Println(<-ch) //load this channel into this print statement
}
