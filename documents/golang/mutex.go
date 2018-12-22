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

	// create a mutex
	var mu sync.Mutex

	for i := 0; i < gs; i++ {
		go func(){
			// "locks" (or "checks out" (like in a Library)) these variables so no other routines can
			// modify them while in use. Everyting from here down to the unlock is considered "locked"
			mu.Lock()
			v := counter
			// This is like time.Sleep(time.Second) - this is cleaner. `Gosched` yeilds the processor to allow
			// other Goroutines to run...this is helpful if you don't have a lot of processors
			runtime.Gosched()
			v++
			counter = v
			// once I'm done with the variables I unlock them (or "check them in") for other goroutines
			// to use and modify.
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Go Routines END:", runtime.NumGoroutine())
	fmt.Println("My Count", counter)
}
