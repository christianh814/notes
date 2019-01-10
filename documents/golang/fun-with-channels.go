package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	go func(){
		c <- 42
		close(c)
	}()

	v, ok := <-c

	// This will print 42 and true because there is a value in "c"
	fmt.Println(v, ok)

	// Let's try to load it again
	v, ok = <-c

	// This will print 0 and false because there is no value (zero value default) in c anymore 
	// because we closed it in the func()...and since there is no value it returns false
	fmt.Println(v, ok)

}
