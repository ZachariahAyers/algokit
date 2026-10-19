// Package combinatorics enumerates products, combinations, subsets and
// permutations.
//
// Enumerations are returned as iterators ([iter.Seq]) because their output
// grows exponentially: a caller can stop early or stream instead of holding
// every result. Each yielded slice is freshly allocated and safe to keep.
package combinatorics

import (
	"cmp"
	"iter"
	"slices"
)

// CartesianProduct yields every tuple that takes one element from each
// set, in lexicographic order of set positions. With no sets it yields one
// empty tuple; if any set is empty it yields nothing.
//
// Time O(k·Π|set|) for k sets. Space O(k) besides the yielded tuples.
func CartesianProduct[T any](sets ...[]T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		for _, s := range sets {
			if len(s) == 0 {
				return
			}
		}
		idx := make([]int, len(sets))
		for {
			tuple := make([]T, len(sets))
			for i, s := range sets {
				tuple[i] = s[idx[i]]
			}
			if !yield(tuple) {
				return
			}
			// Advance the odometer from the last position.
			i := len(sets) - 1
			for ; i >= 0; i-- {
				idx[i]++
				if idx[i] < len(sets[i]) {
					break
				}
				idx[i] = 0
			}
			if i < 0 {
				return
			}
		}
	}
}

// Combinations yields every k-element subset of items, preserving the
// order of items within each subset, in lexicographic order of positions.
// It yields nothing when k < 0 or k > len(items), and one empty subset
// when k == 0.
//
// Time O(k·C(n, k)). Space O(k) besides the yielded subsets.
func Combinations[T any](items []T, k int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(items)
		if k < 0 || k > n {
			return
		}
		idx := make([]int, k)
		for i := range idx {
			idx[i] = i
		}
		for {
			combo := make([]T, k)
			for i, j := range idx {
				combo[i] = items[j]
			}
			if !yield(combo) {
				return
			}
			// Find the rightmost index that can still move right.
			i := k - 1
			for i >= 0 && idx[i] == n-k+i {
				i--
			}
			if i < 0 {
				return
			}
			idx[i]++
			for j := i + 1; j < k; j++ {
				idx[j] = idx[j-1] + 1
			}
		}
	}
}

// PowerSet yields all 2^n subsets of items, ordered by size and then
// lexicographically by position, starting with the empty set.
//
// Time O(n·2^n). Space O(n) besides the yielded subsets.
func PowerSet[T any](items []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		for k := 0; k <= len(items); k++ {
			for combo := range Combinations(items, k) {
				if !yield(combo) {
					return
				}
			}
		}
	}
}

// Permutations yields the distinct permutations of items in ascending
// lexicographic order, using the next-permutation algorithm on a sorted
// copy. Repeated values produce no duplicate permutations, so "aab" yields
// 3 results, not 6. items is not modified.
//
// Time O(n) amortized per permutation, n!/(Π mult!) permutations. Space
// O(n) besides the yielded slices.
func Permutations[T cmp.Ordered](items []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		p := slices.Sorted(slices.Values(items))
		for {
			if !yield(slices.Clone(p)) {
				return
			}
			if !nextPermutation(p) {
				return
			}
		}
	}
}

// nextPermutation rearranges p into the next lexicographic permutation and
// reports false when p was already the last one.
func nextPermutation[T cmp.Ordered](p []T) bool {
	i := len(p) - 2
	for i >= 0 && p[i] >= p[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(p) - 1
	for p[j] <= p[i] {
		j--
	}
	p[i], p[j] = p[j], p[i]
	slices.Reverse(p[i+1:])
	return true
}

// CombinationSum returns every multiset of candidates that sums to target,
// where each candidate may be used any number of times. Each result is in
// ascending order and the results are in lexicographic order. Candidates
// that are not positive are ignored, and duplicate candidates count once.
//
// Time exponential in target / min(candidates) in the worst case; the
// search prunes as soon as a partial sum exceeds target.
func CombinationSum(candidates []int, target int) [][]int {
	var pos []int
	for _, c := range candidates {
		if c > 0 {
			pos = append(pos, c)
		}
	}
	slices.Sort(pos)
	pos = slices.Compact(pos)

	var out [][]int
	var path []int
	var walk func(start, remaining int)
	walk = func(start, remaining int) {
		if remaining == 0 {
			out = append(out, slices.Clone(path))
			return
		}
		for i := start; i < len(pos) && pos[i] <= remaining; i++ {
			path = append(path, pos[i])
			walk(i, remaining-pos[i])
			path = path[:len(path)-1]
		}
	}
	if target > 0 {
		walk(0, target)
	}
	return out
}
