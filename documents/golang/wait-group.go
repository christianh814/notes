package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// create a waitgroup and put it into a variable
var wg sync.WaitGroup

func main() {
	fmt.Println("OS:\t\t", runtime.GOOS)
	fmt.Println("ARCH:\t\t", runtime.GOARCH)
	fmt.Println("CPUs:\t\t", runtime.NumCPU())
	fmt.Println("Goroutines:\t", runtime.NumGoroutine())

	// create a waitgroup and add 1 (wait for one thing)
	wg.Add(1)

	//putting "go" in front of a function launches it in it's own go routine
	go foo()
	bar()

	// print out how many routines are running
	fmt.Println("----------------------------------")
	fmt.Println("Goroutines:\t", runtime.NumGoroutine())

	// Wait for the process to singnal that it's done
	fmt.Println("Waiting for process foo() to be done")
	wg.Wait()
}

func foo() {
	for i := 0; i < 10; i++ {
		fmt.Println("foo:", i)
		// Sleeping for (e/a)ffect
		time.Sleep(1000 * time.Millisecond)
	}
	// This signals that the process is done and for the wg.Wait() to stop waiting
	wg.Done()
}

func bar() {
	for i := 0; i < 10; i++ {
		fmt.Println("bar:", i)
	}
}
