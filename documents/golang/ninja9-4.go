package main

import (
	"fmt"
	"sync"
)

func main() {
	counter := 0
	const gs = 100

	var wg sync.WaitGroup
	wg.Add(gs)

	var mu sync.Mutex

	for i := 0; i < gs; i++ {
		go func() {
			mu.Lock()
			v := counter
			v++
			counter = v
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}
