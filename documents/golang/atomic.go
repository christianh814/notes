package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("CPUs:", runtime.NumCPU())
	fmt.Println("Go Routines:", runtime.NumGoroutine())

	var counter int64 //usually when you see int64 think "atomic"
	const gs = 100

	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func(){
			// atomic takes an address to a conunter and the delta (or how much you want to incrament it by...you can use
			// negatives to decrement)
			atomic.AddInt64(&counter, 1)
			// This is like time.Sleep(time.Second) - this is cleaner. `Gosched` yeilds the processor to allow
			// other Goroutines to run...this is helpful if you don't have a lot of processors
			runtime.Gosched()

			//Let's print it. You need to load the atomic counter. It also takes the address
			fmt.Println(atomic.LoadInt64(&counter))
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Go Routines END:", runtime.NumGoroutine())
	fmt.Println("My Count", counter)
}
