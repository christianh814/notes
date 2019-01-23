package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	//number of bytes written, and the error
	n, err := fmt.Println("hello")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)

	//
	var answer1 string

	// "Scan" takes in input from stdin (like "read" in bash)
	fmt.Print("Name: ")
	_, err = fmt.Scan(&answer1)
	if err != nil {
		panic(err)
	}

	fmt.Println(answer1)

	// Create a file in this path
	f, err := os.Create("/tmp/names.txt")
	// check to see if the file was actually created
	if err != nil {
		fmt.Println(err)
		return // exit if you find an error
	}
	// Defer closing of the file until this function (main()) finishes
	defer f.Close()
	// NewReader returns a new Reader reading from "hello world"
	r := strings.NewReader("James Bond")
	// This writes the text into the file
	io.Copy(f, r)
}
