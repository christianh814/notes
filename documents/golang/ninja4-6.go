package main

import "fmt"

func main() {
	x := make([]string, 4, 50)
	x = append(x, "Alabama", "Alaska", "Arizona", "Arkansas", "California")
	for i := 0; 0 < len(x); i++ {
		fmt.Println(x[i])
	}
}
