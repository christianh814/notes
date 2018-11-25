package main

import "fmt"

var favSport string

func main() {
	favSport = "surfing"

	switch favSport {
	case "canoeing":
		fmt.Println("This will no print")
	case "surfing":
		fmt.Println("print")
	}

}
