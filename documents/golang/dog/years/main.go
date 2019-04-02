// Package years converts dog years into the specified conversion value
package years

import "fmt"

// ToHuman converts dog years into human years
func ToHuman(x int64) (int64, error) {
	if x < 1 {
		return 0, fmt.Errorf("Cannot multiply 0")
	} else {
		return x * 7, nil
	}
}
