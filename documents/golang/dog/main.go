package main

import (
	"fmt"
	"github.com/notes/documents/golang/dog/years"
)

func main() {
	x, err := years.ToHuman(3)
	if err != nil {
	  fmt.Println(err)
	} else {
	  fmt.Printf("In human years I am %v years old\n", x)
	}
}
