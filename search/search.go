// Package search finds a target in a slice.
//
// Every function returns the index of the target and whether it was found.
// The sorted-input searches return the lowest matching index, and when the
// target is absent they return the index where it would be inserted, the
// same contract as [slices.BinarySearch].
package search

import (
	"cmp"
	"math"

	"github.com/ZachariahAyers/algokit/v2/internal/num"
)

// Linear scans s from the front and returns the first index holding target.
// When target is absent it returns len(s), false.
//
// Time O(n). Space O(1). s need not be sorted.
func Linear[T comparable](s []T, target T) (int, bool) {
	for i, v := range s {
		if v == target {
			return i, true
		}
	}
	return len(s), false
}

// Binary searches ascending s by halving the candidate range.
//
// Time O(log n). Space O(1).
func Binary[T cmp.Ordered](s []T, target T) (int, bool) {
	lo, hi := 0, len(s)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if s[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo, lo < len(s) && s[lo] == target
}

// Jump searches ascending s by stepping √n elements at a time, then
// scanning back within the block that must hold target.
//
// Time O(√n). Space O(1). Useful when stepping backwards is expensive.
func Jump[T cmp.Ordered](s []T, target T) (int, bool) {
	n := len(s)
	step := max(1, int(math.Sqrt(float64(n))))
	prev := 0
	for prev < n && s[min(prev+step, n)-1] < target {
		prev += step
	}
	for i := prev; i < min(prev+step, n); i++ {
		if s[i] >= target {
			return i, s[i] == target
		}
	}
	return n, false
}

// Interpolation searches ascending integers by estimating where target
// sits from the values at the ends of the range.
//
// Time O(log log n) on uniformly distributed keys, O(n) worst case.
// Space O(1).
func Interpolation[T num.Integer](s []T, target T) (int, bool) {
	lo, hi := 0, len(s)-1
	for lo <= hi && target >= s[lo] && target <= s[hi] {
		if s[hi] == s[lo] {
			break
		}
		// float64 keeps the estimate from overflowing for wide ranges.
		frac := (float64(target) - float64(s[lo])) / (float64(s[hi]) - float64(s[lo]))
		pos := lo + int(frac*float64(hi-lo))
		pos = min(max(pos, lo), hi)
		switch {
		case s[pos] < target:
			lo = pos + 1
		case s[pos] > target:
			hi = pos - 1
		default:
			// Walk left to the lowest index with this value.
			for pos > lo && s[pos-1] == target {
				pos--
			}
			return pos, true
		}
	}
	// The probe has narrowed s[lo:hi+1]; finish with a bounded binary search
	// so the insertion point is exact.
	i, found := Binary(s[lo:max(hi+1, lo)], target)
	return lo + i, found
}
