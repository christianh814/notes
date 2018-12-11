package main

import (
	"fmt"
	// You need to "go get" this on the cli like so: go get -u golang.org/x/crypto/bcrypt
	"golang.org/x/crypto/bcrypt"
)

func main() {
	p := `password`
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
	fmt.Println(bs)

	givenPw := `passwordBAD`
	berr := bcrypt.CompareHashAndPassword(bs, []byte(givenPw))
	if berr != nil {
		fmt.Println("Bad Password")
	}
}
