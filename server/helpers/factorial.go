package helpers

import "github.com/JohnCGriffin/overflow"

func Factorial(n int64) (int64, bool) {
	// the factorial function is not defined below n = 1
	if n < 1 {
		return 0, false
	}
	var total = n
	var ok bool
	for i := n - 1; i > 0; i-- {
		total, ok = overflow.Mul64(total, i)
		if !ok {
			return 0, ok
		}
	}
	return total, ok
}