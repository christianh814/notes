package main

import (
	"errors"
	"log"
)


func main() {
	// throw away the value and just get the error
	_, err := sqrt(-10)
	// log if there is an error
	if err != nil {
		log.Fatalln(err)
	}
}

func sqrt(f float64) (float64, error) {
	// if `f` i lower than 0 then it's a negative and can't square root it
	if f < 0 {
		return 0, errors.New("nograte math: square root of negative number")
	}
	return 42, nil
}
