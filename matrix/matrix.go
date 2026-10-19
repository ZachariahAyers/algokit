// Package matrix implements in-place operations on square matrices stored
// as [][]T. Aggregates over a whole matrix are array.Flatten followed by
// the matching function in package array.
package matrix

import (
	"errors"
	"slices"
)

// ErrNotSquare reports a matrix whose rows do not all have len(m) columns.
var ErrNotSquare = errors.New("matrix: not square")

// Rotate turns the square matrix m 90 degrees clockwise in place, by
// transposing and then reversing each row.
//
// Time O(n²). Space O(1). Returns [ErrNotSquare] and leaves m unchanged
// if any row length differs from len(m).
func Rotate[T any](m [][]T) error {
	n := len(m)
	for _, row := range m {
		if len(row) != n {
			return ErrNotSquare
		}
	}
	for i := range n {
		for j := i + 1; j < n; j++ {
			m[i][j], m[j][i] = m[j][i], m[i][j]
		}
	}
	for _, row := range m {
		slices.Reverse(row)
	}
	return nil
}
