package numeric

import "math/bits"

// Pow returns base raised to exp by binary exponentiation.
//
// Time O(log exp). Space O(1). Returns [ErrOverflow] if the result does not
// fit in uint64. Pow(0, 0) is 1.
func Pow(base, exp uint64) (uint64, error) {
	result := uint64(1)
	for exp > 0 {
		if exp&1 == 1 {
			hi, lo := bits.Mul64(result, base)
			if hi != 0 {
				return 0, ErrOverflow
			}
			result = lo
		}
		exp >>= 1
		if exp > 0 {
			hi, lo := bits.Mul64(base, base)
			if hi != 0 {
				return 0, ErrOverflow
			}
			base = lo
		}
	}
	return result, nil
}

// IsPowerOf reports whether n == base^k for some k >= 0.
//
// Time O(log n / log base). Space O(1). Base 0 matches 0 and 1; base 1
// matches only 1.
func IsPowerOf(n, base uint64) bool {
	switch base {
	case 0:
		return n <= 1
	case 1:
		return n == 1
	}
	if n == 0 {
		return false
	}
	for n%base == 0 {
		n /= base
	}
	return n == 1
}

// Sqrt returns the integer square root of n: the largest r with r*r <= n.
// It uses Newton's method, which converges quadratically from above.
//
// Time O(log log n) iterations after the bit-length seed. Space O(1).
func Sqrt(n uint64) uint64 {
	if n < 2 {
		return n
	}
	// Seed with a power of two at or above √n.
	x := uint64(1) << ((bits.Len64(n) + 1) / 2)
	for {
		y := (x + n/x) / 2
		if y >= x {
			return x
		}
		x = y
	}
}
