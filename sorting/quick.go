package sorting

import "cmp"

// Quick sorts s with quicksort: median-of-three pivot, three-way partition
// so runs of equal keys cost linear time, and recursion on the smaller side
// only so the stack stays O(log n).
//
// Time O(n log n) expected, O(n²) worst case. Space O(log n). Not stable.
func Quick[T cmp.Ordered](s []T) {
	for len(s) > 1 {
		lt, gt := partition(s)
		if lt < len(s)-gt {
			Quick(s[:lt])
			s = s[gt:]
		} else {
			Quick(s[gt:])
			s = s[:lt]
		}
	}
}

// partition rearranges s into < pivot, == pivot, > pivot and returns the
// bounds of the middle band: s[lt:gt] all equal the pivot.
func partition[T cmp.Ordered](s []T) (lt, gt int) {
	pivot := medianOfThree(s[0], s[len(s)/2], s[len(s)-1])
	lt, i, gt := 0, 0, len(s)
	for i < gt {
		switch {
		case s[i] < pivot:
			s[lt], s[i] = s[i], s[lt]
			lt++
			i++
		case s[i] > pivot:
			gt--
			s[i], s[gt] = s[gt], s[i]
		default:
			i++
		}
	}
	return lt, gt
}

func medianOfThree[T cmp.Ordered](a, b, c T) T {
	if a > b {
		a, b = b, a
	}
	if b > c {
		b = c
	}
	return max(a, b)
}
