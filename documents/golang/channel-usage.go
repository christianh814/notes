package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

	//send
	go send(c)

	//recv
	recv(c)

	fmt.Println("Exiting")
}

func send(c chan<- int) {
	c <- 42
}

func recv(c <-chan int) {
	fmt.Println(<-c)
}
