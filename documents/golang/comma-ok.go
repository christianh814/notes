package main

import (
	"fmt"
)

func main() {
	eve := make(chan int)
	odd := make(chan int)
	quit := make(chan bool)

	// send
	go send(eve, odd, quit)

	// recv
	recv(eve, odd, quit)

	fmt.Println("Done running")
}

func send(e, o chan<- int, q chan<- bool) {
	for i := 0; i < 100; i++ {
		if (i % 2) == 0 {
			e <- i
		} else {
			o <- i
		}
	}
	close(q)
}

func recv(e, o <-chan int, q <-chan bool) {
	for {
		select {
		case v := <-e:
			fmt.Println("From the even:", v)
		case v := <-o:
			fmt.Println("From the odd:", v)
		case i, ok := <-q:
			if !ok {
				fmt.Println("From OK comma", i)
				return
			} else {
				fmt.Println("From OK comma", i)
			}
		}
	}
}
