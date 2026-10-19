package numeric

import "math"

// ReverseDigits returns n with its decimal digits reversed, keeping the
// sign: -120 becomes -21.
//
// Time O(d) for d digits. Space O(1). Returns [ErrOverflow] when the
// reversed value does not fit in int64.
func ReverseDigits(n int64) (int64, error) {
	var r int64
	for n != 0 {
		d := n % 10
		if r > math.MaxInt64/10 || r < math.MinInt64/10 ||
			r == math.MaxInt64/10 && d > 7 || r == math.MinInt64/10 && d < -8 {
			return 0, ErrOverflow
		}
		r = r*10 + d
		n /= 10
	}
	return r, nil
}

// SameDigits reports whether a and b use the same multiset of decimal
// digits, ignoring sign: 1123 and 3121 match, 12 and 122 do not.
//
// Time O(d). Space O(1).
func SameDigits(a, b int64) bool {
	var counts [10]int
	for ; a != 0; a /= 10 {
		counts[abs(a%10)]++
	}
	for ; b != 0; b /= 10 {
		counts[abs(b%10)]--
	}
	return counts == [10]int{}
}

func abs(d int64) int64 {
	if d < 0 {
		return -d
	}
	return d
}
