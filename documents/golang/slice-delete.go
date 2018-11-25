package main

import "fmt"

func main() {
	x := []int{4, 5, 42, 7, 8, 12, 99}
	fmt.Println(x)
	// Delete 7 and 8. give me from start to pos 2 and append pos 5 onwards
	x = append(x[:3], x[5:]...)
	fmt.Println(x)
}
