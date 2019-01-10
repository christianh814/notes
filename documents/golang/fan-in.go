package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create your channels
	even := make(chan int)
	odd := make(chan int)
	fanin := make(chan int)

	// run "send in" and "recv" both in their own respective go routines
	go send(even, odd)
	go receive(even, odd, fanin)

	for v := range fanin {
		fmt.Println(v)
	}

	fmt.Println("about to exit")
}

// send channel
func send(even, odd chan<- int) {
	// This stores even numbers in the "even" channel and
	// the odd numbers in the "odd" channel
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			even <- i
		} else {
			odd <- i
		}
	}
	// Once done we close our channels
	close(even)
	close(odd)
}

// receive channel
func receive(even, odd <-chan int, fanin chan<- int) {
	// Create a waitgroup
	var wg sync.WaitGroup
	// We are using 2 go routines here so we need to create a waitgroup of 2
	wg.Add(2)

	// go routine that loads the even numbers from the "even" channel into the "fanin" channel
	go func() {
		for v := range even {
			fanin <- v
		}
		// Flag the waitgroup that this goroutine is done
		wg.Done()
	}()

	// go routine that loads the odd numbers from the "odd" channel into the "fanin" channel
	go func() {
		for v := range odd {
			fanin <- v
		}
		// Flag the waitgroup that this goroutine is done
		wg.Done()
	}()

	// wait for the goroutines to finish
	wg.Wait()

	// close the "fanin" channel
	close(fanin)
}

