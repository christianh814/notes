package main

import "fmt"

func main() {
	x := []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51}
	fmt.Println(x)
	//sclice it up so you get what you're looking for
	x = append(x[:3], x[6:]...)
	fmt.Println(x)
}
