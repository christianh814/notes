package main

import (
	"fmt"
	"os"
	"io"
)

func main() {
	fmt.Println("Hello World")
	fmt.Fprintln(os.Stdout, "Hello worlds")
	io.WriteString(os.Stdout, "hello io\n")
}
