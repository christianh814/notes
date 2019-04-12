package main

import (
	"testing"
)

func TestMySum(t *testing.T) {

	type test struct {
		data   []int
		answer int
	}

	tests := []test{
		test{[]int{21, 21}, 42},
		test{[]int{1, 1}, 2},
		test{[]int{2, 2}, 4},
		test{[]int{11, 11}, 22},
		test{[]int{5, 5}, 10},
		test{[]int{3, 3}, 6},
	}

	for _, v := range tests {
		x := mySum(v.data...)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}

	}

}
