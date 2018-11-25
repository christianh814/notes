package main

import "fmt"

const (
	// _ "throws away" the first iota (i.e 0)
	_ = iota
	// take one binary and shift it over 10 times
	kilo = 1 << (iota * 10)
	// take one binary and shift it over 20 times
	mega = 1 << (iota * 10)
	// take one binary and shift it over 30 times
	giga = 1 << (iota * 10)
)

func main() {
	kb := 1024
	mb := 1024 * kb
	gb := 1024 * mb

	fmt.Printf("%d\t\t%b\n", kb, kb)
	fmt.Printf("%d\t\t%b\n", mb, mb)
	fmt.Printf("%d\t%b\n", gb, gb)

	fmt.Printf("%d\t\t%b\n", kilo, kilo)
	fmt.Printf("%d\t\t%b\n", mega, mega)
	fmt.Printf("%d\t%b\n", giga, giga)
}
