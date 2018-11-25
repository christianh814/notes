package main

import "fmt"

func main() {
	a := (42 == 42)
	b := (42 <= 100)
	c := (42 >= 31)
	d := (42 != 100)
	e := (1 < 2)
	f := (5 > 1)

	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\n", a, b, c, d, e, f)
}
