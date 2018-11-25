package main

import "fmt"

type person struct {
	fn string
	ln string
}

type secretAgent struct {
	person
	ltk bool
}

//  attach function speak() to the type secretAgent. So it's available to any secretAgent
func (s secretAgent) speak() {
	fmt.Print("I am ", s.ln, ", ", s.fn, " ", s.ln)
	if s.ltk {
		fmt.Println("...I have a license to kill")
	}
}
//
func main() {
	sa1 := secretAgent{
		person: person{
			"James",
			"Bond",
		},
		ltk: true,
	}
	fmt.Println(sa1)
	sa1.speak()
}
