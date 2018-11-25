package main

import "fmt"

func main() {
	jb := []string{"sun", "mon", "tue"}
	fmt.Println(jb)
	//
	mp := []string{"wed", "thr", "fr"}
	fmt.Println(mp)
	// a slice of a slice of string
	xp := [][]string{jb, mp}
	fmt.Println(xp)
}
