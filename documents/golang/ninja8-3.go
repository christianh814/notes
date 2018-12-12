package main

import (
	"fmt"
	"encoding/json"
	"os"
)

type user struct {
	First string
	Last string
	Age int
}

func main() {
	u1 := user{
		First: "Christian",
		Last: "Hernandez",
		Age: 36,
	}
	u2 := user{
		First: "Mark",
		Last: "Anthony",
		Age: 45,
	}
	u3 := user{
		First: "Martha",
		Last: "Smith",
		Age: 24,
	}
	users := []user{u1, u2, u3}
	// this turns the `users` data structure (which is a slice of type user...which is a struct) into a json format. Pipe this to `python -m json.tool` to see
	err := json.NewEncoder(os.Stdout).Encode(users)
	if err != nil {
		fmt.Println(err)
	}
}
