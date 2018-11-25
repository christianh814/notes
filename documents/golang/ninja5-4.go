package main

import "fmt"

func main() {
	s := struct {
		first string
		last string
		hands map[string]string
	}{
		first: "Christian",
		last: "Hernandez",
		hands: map[string]string{
			"left": "fine",
			"right": "swollen",
		},
	}
	fmt.Println("It seems that", s.first, s.last, "is", s.hands["left"])
}
