package main

import (
	"fmt"
	"encoding/json"
)

// In order for these to be exported, they need to be capital
type user struct{
	First string
	Age int
}

func main() {
	u1 := user{
		First: "James",
		Age: 47,
	}
	u2 := user{
		First: "Ricky",
		Age: 22,
	}
	u3 := user{
		First: "Peter",
		Age: 77,
	}
	users := []user{u1, u2, u3}
	//
	ue, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(ue))
}
