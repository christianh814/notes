package main

import "fmt"

func foo() string {
	return "everyone has the func"
}

func bar(f func() string) {
	w := f()
	fmt.Println(w)
}

func main() {
	bar(foo())
}
