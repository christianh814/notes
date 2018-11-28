package main

import (
	"fmt"
	"math"
)

type square struct {
	l float64
	w float64
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (s square) area() float64 {
	return s.l * s.w
}

type shape interface {
	area() float64
}

func info(s shape) {
	fmt.Println(s.area())
}

func main() {
	c := circle{
		radius: 12.345,
	}
	s := square{
		l: 3.5,
		w: 7.8,
	}
	info(c)
	info(s)
}
