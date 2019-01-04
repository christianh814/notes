package main

import "fmt"

func main() {
	p := struct {
		first string
		last  string
		age   int
	}{
		first: "Christian",
		last:  "Hernandez",
		age:   36,
	}
	fmt.Println(p.first)
}
