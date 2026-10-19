package numeric

import (
	"math/bits"
	"strconv"
)

// Factorial returns n!.
//
// Time O(n). Space O(1). Returns [ErrOverflow] for n > 20, the largest
// factorial that fits in uint64.
func Factorial(n uint) (uint64, error) {
	f := uint64(1)
	for i := uint64(2); i <= uint64(n); i++ {
		hi, lo := bits.Mul64(f, i)
		if hi != 0 {
			return 0, ErrOverflow
		}
		f = lo
	}
	return f, nil
}

// Fibonacci returns F(n), with F(0) = 0 and F(1) = 1.
//
// Time O(n). Space O(1). Returns [ErrOverflow] for n > 93.
func Fibonacci(n uint) (uint64, error) {
	if n > 93 {
		return 0, ErrOverflow
	}
	a, b := uint64(0), uint64(1)
	for range n {
		// b may wrap on the final step; only a is returned.
		a, b = b, a+b
	}
	return a, nil
}

// CollatzSteps returns how many steps the Collatz map (n/2 if even, 3n+1
// if odd) takes to reach 1 from n.
//
// Time proportional to the step count, which is unproven to be finite in
// general but is for every uint64 tested. Returns [ErrDomain] for n == 0
// and [ErrOverflow] if 3n+1 leaves uint64.
func CollatzSteps(n uint64) (int, error) {
	if n == 0 {
		return 0, ErrDomain
	}
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			if n > (^uint64(0)-1)/3 {
				return 0, ErrOverflow
			}
			n = 3*n + 1
		}
		steps++
	}
	return steps, nil
}

// FizzBuzz returns the FizzBuzz sequence for 1..n: "Fizz" for multiples of
// 3, "Buzz" for multiples of 5, "FizzBuzz" for both, otherwise the number.
//
// Time O(n). Space O(n).
func FizzBuzz(n int) []string {
	out := make([]string, 0, max(n, 0))
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			out = append(out, "FizzBuzz")
		case i%3 == 0:
			out = append(out, "Fizz")
		case i%5 == 0:
			out = append(out, "Buzz")
		default:
			out = append(out, strconv.Itoa(i))
		}
	}
	return out
}
