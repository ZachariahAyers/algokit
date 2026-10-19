package sorting

import (
	"slices"

	"github.com/ZachariahAyers/algokit/v2/internal/num"
)

// Counting sorts integers with counting sort. Negative values are allowed.
//
// Time O(n + k) and space O(k), where k = max(s) - min(s) + 1. Only suitable
// when k is small next to n. Panics if k exceeds [MaxCountingRange].
func Counting[T num.Integer](s []T) {
	if len(s) < 2 {
		return
	}
	lo, hi := slices.Min(s), slices.Max(s)
	// Unsigned subtraction wraps, so the span is right for every signed width.
	span := uint64(hi) - uint64(lo)
	if span >= MaxCountingRange {
		panic("sorting: Counting value range exceeds MaxCountingRange")
	}
	counts := make([]int, span+1)
	for _, v := range s {
		counts[uint64(v)-uint64(lo)]++
	}
	i := 0
	for off, c := range counts {
		for ; c > 0; c-- {
			s[i] = lo + T(off)
			i++
		}
	}
}

// MaxCountingRange caps the number of counters [Counting] will allocate.
const MaxCountingRange = 1 << 26

// Radix sorts unsigned integers with least-significant-digit radix sort,
// one byte per pass.
//
// Time O(w·n) and space O(n), where w is the byte width of T. Stable.
func Radix[T num.Unsigned](s []T) {
	if len(s) < 2 {
		return
	}
	hi := slices.Max(s)
	buf := make([]T, len(s))
	src, dst := s, buf
	for shift := uint(0); shift < 64 && hi>>shift > 0; shift += 8 {
		var counts [257]int
		for _, v := range src {
			counts[int(byte(v>>shift))+1]++
		}
		for b := 1; b < len(counts); b++ {
			counts[b] += counts[b-1]
		}
		for _, v := range src {
			d := byte(v >> shift)
			dst[counts[d]] = v
			counts[d]++
		}
		src, dst = dst, src
	}
	if &src[0] != &s[0] {
		copy(s, src)
	}
}
