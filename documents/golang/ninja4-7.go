package main

import "fmt"

func main() {
	a := []string{"James", "Hello James"}
	b := []string{"Miss", "Pennyfeather"}
	xm := [][]string{a, b}
	fmt.Println(xm)

	for i, v := range xm {
		fmt.Println(i, v)
	}
}
