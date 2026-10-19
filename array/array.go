// Package array implements whole-slice algorithms: aggregates, duplicate
// detection, rotation, shuffling and compaction.
//
// Minimum and maximum are left to [slices.Min] and [slices.Max].
package array

import (
	"errors"
	"math/rand/v2"
	"slices"

	"github.com/ZachariahAyers/algokit/v2/internal/num"
)

// ErrEmpty reports an aggregate that has no value for an empty slice.
var ErrEmpty = errors.New("array: empty slice")

// Flatten concatenates nested slices into one new slice, in order.
//
// Time O(n) for n total elements. Space O(n).
func Flatten[T any](nested [][]T) []T {
	n := 0
	for _, s := range nested {
		n += len(s)
	}
	out := make([]T, 0, n)
	for _, s := range nested {
		out = append(out, s...)
	}
	return out
}

// Sum returns the sum of s, 0 when s is empty. Integer sums wrap on
// overflow, as Go integer arithmetic does.
//
// Time O(n). Space O(1).
func Sum[T num.Real](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

// Product returns the product of s, 1 when s is empty. Integer products
// wrap on overflow.
//
// Time O(n). Space O(1).
func Product[T num.Real](s []T) T {
	total := T(1)
	for _, v := range s {
		total *= v
	}
	return total
}

// Mean returns the arithmetic mean of s.
//
// Time O(n). Space O(1). Returns [ErrEmpty] for an empty slice.
func Mean[T num.Real](s []T) (float64, error) {
	if len(s) == 0 {
		return 0, ErrEmpty
	}
	// Accumulate in float64 so integer inputs cannot overflow the sum.
	var total float64
	for _, v := range s {
		total += float64(v)
	}
	return total / float64(len(s)), nil
}

// Median returns the middle value of s, or the mean of the two middle
// values when len(s) is even. s is not modified.
//
// Time O(n log n) for the sort of a copy. Space O(n). Returns [ErrEmpty]
// for an empty slice.
func Median[T num.Real](s []T) (float64, error) {
	if len(s) == 0 {
		return 0, ErrEmpty
	}
	sorted := slices.Clone(s)
	slices.Sort(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[mid]), nil
	}
	return (float64(sorted[mid-1]) + float64(sorted[mid])) / 2, nil
}

// Duplicates returns each value that occurs more than once in s, once,
// in order of its second occurrence.
//
// Time O(n). Space O(n).
func Duplicates[T comparable](s []T) []T {
	seen := make(map[T]int, len(s))
	var out []T
	for _, v := range s {
		seen[v]++
		if seen[v] == 2 {
			out = append(out, v)
		}
	}
	return out
}

// IsUnique reports whether every value in s is distinct.
//
// Time O(n). Space O(n).
func IsUnique[T comparable](s []T) bool {
	seen := make(map[T]struct{}, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			return false
		}
		seen[v] = struct{}{}
	}
	return true
}

// Rotate rotates s left by k places in place, so s[k] becomes s[0].
// A negative k rotates right. It uses three reversals.
//
// Time O(n). Space O(1).
func Rotate[T any](s []T, k int) {
	n := len(s)
	if n == 0 {
		return
	}
	k = ((k % n) + n) % n
	slices.Reverse(s[:k])
	slices.Reverse(s[k:])
	slices.Reverse(s)
}

// Shuffle permutes s uniformly at random in place with the Fisher-Yates
// algorithm, drawing from r so callers control the seed.
//
// Time O(n). Space O(1).
func Shuffle[T any](s []T, r *rand.Rand) {
	for i := len(s) - 1; i > 0; i-- {
		j := r.IntN(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// MoveZeroes moves every zero value of s to the end in place, keeping the
// other elements in their original order.
//
// Time O(n). Space O(1).
func MoveZeroes[T comparable](s []T) {
	var zero T
	w := 0
	for _, v := range s {
		if v != zero {
			s[w] = v
			w++
		}
	}
	clear(s[w:])
}
