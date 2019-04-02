package main

import (
	"fmt"
	// this path is relative to $GOPATH
	"github.com/notes/documents/golang/documentation_go/greaterthan"
)

func main() {
	a := 10
	b := 2
	if greaterthan.Gt(a, b) {
		fmt.Printf("Number %v is greater than %v\n", a, b)
	}
}
