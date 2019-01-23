package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	for j := 0; j < 10; j++ {
		go func() {
			for i := 0; i <= 9; i++ {
				c <- i
			}
		}()
	}
	//
	for k := 0; k<100; k++ {
		fmt.Println(<-c)
	}
}
