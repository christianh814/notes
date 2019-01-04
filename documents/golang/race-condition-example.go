package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("CPUs:", runtime.NumCPU())
	fmt.Println("Go Routines:", runtime.NumGoroutine())

	counter := 0
	const gs = 100

	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			v := counter
			// This is like time.Sleep(time.Second) - this is cleaner. `Gosched` yeilds the processor to allow
			// other Goroutines to run...this is helpful if you don't have a lot of processors
			runtime.Gosched()
			v++
			counter = v
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Go Routines END:", runtime.NumGoroutine())
	fmt.Println("My Count", counter)
}
