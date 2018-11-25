package main

import "fmt"

func main() {
	y := []string{"cat", "dog", "mouse"}
	// I can just print the value of the index by doing `for _, v := range y`...basically use _ to "throw away" the index
	for i, v := range y {
		fmt.Println(i, v)
	}
	fmt.Printf("%T\n", y)
}
