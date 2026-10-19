package sorting

import "cmp"

// Bubble sorts s with bubble sort, stopping early once a pass makes no swap.
//
// Time O(n²), O(n) when s is already sorted. Space O(1). Stable.
func Bubble[T cmp.Ordered](s []T) {
	for end := len(s); end > 1; end-- {
		swapped := false
		for i := 1; i < end; i++ {
			if s[i] < s[i-1] {
				s[i], s[i-1] = s[i-1], s[i]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

// Insertion sorts s with insertion sort.
//
// Time O(n²), O(n) when s is already sorted. Space O(1). Stable.
func Insertion[T cmp.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i
		for ; j > 0 && v < s[j-1]; j-- {
			s[j] = s[j-1]
		}
		s[j] = v
	}
}

// Selection sorts s with selection sort.
//
// Time O(n²) on every input. Space O(1). Not stable. Makes at most n-1 swaps.
func Selection[T cmp.Ordered](s []T) {
	for i := range s {
		minIdx := i
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[minIdx] {
				minIdx = j
			}
		}
		s[i], s[minIdx] = s[minIdx], s[i]
	}
}

// Shell sorts s with Shell sort using Ciura's gap sequence, extended by
// factor 2.25 for long inputs.
//
// Time about O(n^1.3) in practice; no tight bound is known for this
// sequence. Space O(1). Not stable.
func Shell[T cmp.Ordered](s []T) {
	gaps := []int{1, 4, 10, 23, 57, 132, 301, 701, 1750}
	for g := gaps[len(gaps)-1]; g < len(s)/2; {
		g = g * 9 / 4
		gaps = append(gaps, g)
	}
	for k := len(gaps) - 1; k >= 0; k-- {
		gap := gaps[k]
		for i := gap; i < len(s); i++ {
			v := s[i]
			j := i
			for ; j >= gap && v < s[j-gap]; j -= gap {
				s[j] = s[j-gap]
			}
			s[j] = v
		}
	}
}
