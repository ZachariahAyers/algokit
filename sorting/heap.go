package sorting

import "cmp"

// Heap sorts s with heapsort: build a max-heap in place, then move the root
// to the end one element at a time.
//
// Time O(n log n) on every input. Space O(1). Not stable.
func Heap[T cmp.Ordered](s []T) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		siftDown(s, i)
	}
	for end := len(s) - 1; end > 0; end-- {
		s[0], s[end] = s[end], s[0]
		siftDown(s[:end], 0)
	}
}

func siftDown[T cmp.Ordered](s []T, i int) {
	for {
		child := 2*i + 1
		if child >= len(s) {
			return
		}
		if r := child + 1; r < len(s) && s[r] > s[child] {
			child = r
		}
		if s[i] >= s[child] {
			return
		}
		s[i], s[child] = s[child], s[i]
		i = child
	}
}
