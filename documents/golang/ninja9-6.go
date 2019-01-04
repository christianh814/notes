package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Println("Arch:", runtime.GOARCH)
}
