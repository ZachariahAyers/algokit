package sorting

import "cmp"

// Merge sorts s with top-down merge sort.
//
// Time O(n log n) on every input. Space O(n) for one scratch buffer. Stable.
func Merge[T cmp.Ordered](s []T) {
	if len(s) < 2 {
		return
	}
	mergeSort(s, make([]T, len(s)))
}

func mergeSort[T cmp.Ordered](s, buf []T) {
	if len(s) < 2 {
		return
	}
	mid := len(s) / 2
	mergeSort(s[:mid], buf[:mid])
	mergeSort(s[mid:], buf[mid:])
	if s[mid-1] <= s[mid] {
		return
	}
	copy(buf, s)
	i, j, k := 0, mid, 0
	for i < mid && j < len(s) {
		// Taking from the left on ties keeps the sort stable.
		if buf[j] < buf[i] {
			s[k] = buf[j]
			j++
		} else {
			s[k] = buf[i]
			i++
		}
		k++
	}
	k += copy(s[k:], buf[i:mid])
	copy(s[k:], buf[j:])
}
