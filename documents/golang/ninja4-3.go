package main

import "fmt"

func main() {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(a)
	fmt.Println(a[1:2])
	fmt.Println(a[3:7])
	fmt.Println(a[:7])
	fmt.Println(a[5:])
	fmt.Println(a[1:4])
}
