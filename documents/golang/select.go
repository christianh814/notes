package main

import (
	"fmt"
)

func main() {
	eve := make(chan int)
	odd := make(chan int)
	quit := make(chan int)

	// send
	go send(eve, odd, quit)

	// recv
	recv(eve, odd, quit)

	fmt.Println("Done running")
}

func send(e, o, q chan<- int) {
	for i := 0; i < 100; i++ {
		if (i % 2) == 0 {
			e <- i
		} else {
			o <- i
		}
	}
	q <- 0
}

func recv(e, o, q <-chan int) {
	for {
		select {
		case v := <-e:
			fmt.Println("From the even:", v)
		case v := <-o:
			fmt.Println("From the odd:", v)
		case v := <-q:
			fmt.Println("From the ZERO:", v)
			return
		}
	}
}
