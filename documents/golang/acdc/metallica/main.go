// Package metallica asks if you're ready to rock
package metallica

// Sum adds an unlimited amount of values of type int
func Sum(xi ...int) int {
	s := 0
	for _, v := range xi {
		s += v
	}
	return s
}
