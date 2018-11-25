package main

import "fmt"

//
type hotdog int

//

type person struct {
	fn string
	ln string
}

type secretAgent struct {
	person
	ltk bool
}

// interfaces - any type that has the method of `speak()` is ALSO a type of `human()` - a VALUE can be more than one TYPE
type human interface {
	speak()
}

// Since both secretAgent and person has the method speak() that means it's also a human...bar() takes in human and I can assertain the type
func bar(h human) {
	switch h.(type) {
	case person:
		fmt.Println("person", h.(person).fn)
	case secretAgent:
		fmt.Println("I cannot say", h.(secretAgent).fn)
	}
	fmt.Println("I am a human ", h)
}

//  attach function speak() to the type secretAgent. So it's available to any secretAgent
func (s secretAgent) speak() {
	fmt.Print("I am ", s.ln, ", ", s.fn, " ", s.ln)
	if s.ltk {
		fmt.Println("...I have a license to kill")
	}
}
func (p person) speak() {
	fmt.Print("I am ", p.ln, ", ", p.fn, " ", p.ln, "...capish?")
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
	//
	p1 := person{
		fn: "Dr",
		ln: "No",
	}

	fmt.Println(p1)
	bar(sa1)
	bar(p1)
	//conversion
	var x hotdog = 32
	fmt.Println(x)
	fmt.Printf("%T\n", x)
	var y int
	y = int(x)
	fmt.Println(y)
	fmt.Printf("%T\n", y)
}

//
// -30-
//
