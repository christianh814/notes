package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	// Open a file in this path
	f, err := os.Open("/tmp/names.txt")
	// check to see if the file was actually opend
	if err != nil {
		fmt.Println(err)
		return // exit if you find an error
	}
	// Defer closing of the file until this function (main()) finishes
	defer f.Close()
	//read the contents of the file (in bytes)
	bs, err := ioutil.ReadAll(f)
	// check if you were successfully able to read it
	if err != nil {
		fmt.Println(err)
		return // exit if you find an error
	}
	//Print the contents...turning the bytes into a string
	fmt.Println(string(bs))
}
