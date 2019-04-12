package main

import (
	"fmt"
	"github.com/christianh814/notes/documents/golang/acdc/metallica"
)

func main() {
	fmt.Println(metallica.Sum(2, 3, 4, 5))
	fmt.Println(metallica.Sum(2, 3, 4, 5, 999))
}
