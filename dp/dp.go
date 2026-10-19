// Package dp implements dynamic-programming algorithms over sequences:
// longest common, increasing and decreasing subsequences, the shortest
// common supersequence, edit distance and the 0/1 knapsack.
package dp

import (
	"cmp"
	"errors"
	"slices"
)

// lcsTable returns t where t[i][j] is the LCS length of a[i:] and b[j:].
func lcsTable[T comparable](a, b []T) [][]int {
	t := make([][]int, len(a)+1)
	for i := range t {
		t[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				t[i][j] = t[i+1][j+1] + 1
			} else {
				t[i][j] = max(t[i+1][j], t[i][j+1])
			}
		}
	}
	return t
}

// LongestCommonSubsequence returns a longest sequence that appears, not
// necessarily contiguously, in both a and b. When several exist it returns
// the one found by preferring to advance a.
//
// Time O(m·n). Space O(m·n).
func LongestCommonSubsequence[T comparable](a, b []T) []T {
	t := lcsTable(a, b)
	out := make([]T, 0, t[0][0])
	for i, j := 0, 0; i < len(a) && j < len(b); {
		switch {
		case a[i] == b[j]:
			out = append(out, a[i])
			i++
			j++
		case t[i+1][j] >= t[i][j+1]:
			i++
		default:
			j++
		}
	}
	return out
}

// ShortestCommonSupersequence returns a shortest sequence that has both a
// and b as subsequences. Its length is len(a) + len(b) - LCS(a, b).
//
// Time O(m·n). Space O(m·n).
func ShortestCommonSupersequence[T comparable](a, b []T) []T {
	t := lcsTable(a, b)
	out := make([]T, 0, len(a)+len(b)-t[0][0])
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			out = append(out, a[i])
			i++
			j++
		case t[i+1][j] >= t[i][j+1]:
			out = append(out, a[i])
			i++
		default:
			out = append(out, b[j])
			j++
		}
	}
	out = append(out, a[i:]...)
	return append(out, b[j:]...)
}

// LongestIncreasingSubsequence returns a longest strictly increasing
// subsequence of s, using patience sorting with binary search. Among
// equally long answers it returns one ending in the smallest possible value.
//
// Time O(n log n). Space O(n).
func LongestIncreasingSubsequence[T cmp.Ordered](s []T) []T {
	return longestMonotone(s, cmp.Compare[T])
}

// LongestDecreasingSubsequence returns a longest strictly decreasing
// subsequence of s.
//
// Time O(n log n). Space O(n).
func LongestDecreasingSubsequence[T cmp.Ordered](s []T) []T {
	return longestMonotone(s, func(a, b T) int { return cmp.Compare(b, a) })
}

// longestMonotone finds a longest subsequence strictly increasing under
// compare.
func longestMonotone[T any](s []T, compare func(a, b T) int) []T {
	// tails[k] is the index in s of the smallest tail of any increasing
	// subsequence of length k+1 seen so far; prev links each index to its
	// predecessor in that subsequence.
	var tails []int
	prev := make([]int, len(s))
	for i, v := range s {
		k, _ := slices.BinarySearchFunc(tails, v, func(idx int, v T) int {
			return compare(s[idx], v)
		})
		if k > 0 {
			prev[i] = tails[k-1]
		} else {
			prev[i] = -1
		}
		if k == len(tails) {
			tails = append(tails, i)
		} else {
			tails[k] = i
		}
	}
	out := make([]T, len(tails))
	if len(tails) == 0 {
		return out
	}
	for k, i := len(tails)-1, tails[len(tails)-1]; k >= 0; k, i = k-1, prev[i] {
		out[k] = s[i]
	}
	return out
}

// EditDistance returns the Levenshtein distance between a and b: the
// fewest single-element insertions, deletions and substitutions that turn
// a into b. For strings, pass []rune(s).
//
// Time O(m·n). Space O(min(m, n)) with two rolling rows.
func EditDistance[T comparable](a, b []T) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i := 1; i <= len(a); i++ {
		diag := row[0]
		row[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			next := min(row[j]+1, row[j-1]+1, diag+cost)
			diag, row[j] = row[j], next
		}
	}
	return row[len(b)]
}

// ErrKnapsackInput reports mismatched or negative knapsack inputs.
var ErrKnapsackInput = errors.New("dp: weights and values must be the same length, weights and capacity non-negative")

// Knapsack solves the 0/1 knapsack problem: choose items, each at most
// once, to maximise total value with total weight at most capacity. It
// returns the best value and the chosen item indexes in ascending order.
//
// Time O(n·W) and space O(n·W) for W = capacity, which is pseudo-polynomial:
// fine for capacities in the thousands, not for 2^40.
func Knapsack(weights, values []int, capacity int) (best int, chosen []int, err error) {
	if len(weights) != len(values) || capacity < 0 {
		return 0, nil, ErrKnapsackInput
	}
	for _, w := range weights {
		if w < 0 {
			return 0, nil, ErrKnapsackInput
		}
	}
	n := len(weights)
	// t[i][c] is the best value using items[i:] with capacity c.
	t := make([][]int, n+1)
	for i := range t {
		t[i] = make([]int, capacity+1)
	}
	for i := n - 1; i >= 0; i-- {
		for c := 0; c <= capacity; c++ {
			t[i][c] = t[i+1][c]
			if weights[i] <= c {
				t[i][c] = max(t[i][c], values[i]+t[i+1][c-weights[i]])
			}
		}
	}
	c := capacity
	for i := range n {
		if t[i][c] != t[i+1][c] {
			chosen = append(chosen, i)
			c -= weights[i]
		}
	}
	return t[0][capacity], chosen, nil
}
