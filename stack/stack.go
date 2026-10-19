// Package stack implements a generic last-in, first-out stack backed by a
// slice. The zero value is an empty stack.
package stack

import "iter"

// Stack is a LIFO stack.
type Stack[T any] struct {
	items []T
}

// Len returns the number of elements. O(1).
func (s *Stack[T]) Len() int { return len(s.items) }

// IsEmpty reports whether the stack has no elements. O(1).
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }

// Push adds v to the top. O(1) amortized.
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop removes and returns the top element. O(1).
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	top := s.items[len(s.items)-1]
	// Clear the slot so a popped pointer does not stay reachable.
	s.items[len(s.items)-1] = zero
	s.items = s.items[:len(s.items)-1]
	return top, true
}

// Peek returns the top element without removing it. O(1).
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Clear removes every element. O(1).
func (s *Stack[T]) Clear() { s.items = nil }

// All yields the elements from top to bottom.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.items) - 1; i >= 0; i-- {
			if !yield(s.items[i]) {
				return
			}
		}
	}
}
