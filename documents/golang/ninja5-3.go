package main

import "fmt"

type vehicle struct{
	doors int
	color string
}

type truck struct {
	vehicle
	fourWheel bool
}

type sedan struct {
	vehicle
	lux bool
}

func main() {
	ford := truck {
		vehicle: vehicle{
			doors: 2,
			color: "green",
		},
		fourWheel: true,
	}
	fmt.Println(ford)
	lexus := sedan {
		vehicle: vehicle{
			doors: 4,
			color: "gold",
		},
		lux: true,
	}
	fmt.Println(lexus)
	//
	fmt.Println(ford.doors)
	fmt.Println(lexus.doors)
	fmt.Println(ford.fourWheel)
	fmt.Println(lexus.lux)
}
