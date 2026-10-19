// Package set implements a generic unordered set on a built-in map.
// The zero value is an empty set.
package set

import (
	"iter"
	"maps"
)

// Set is an unordered collection of distinct values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// Of returns a set holding values.
func Of[T comparable](values ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(values))}
	for _, v := range values {
		s.m[v] = struct{}{}
	}
	return s
}

// Len returns the number of values. O(1).
func (s *Set[T]) Len() int { return len(s.m) }

// Add inserts v and reports whether it was new. O(1) expected.
func (s *Set[T]) Add(v T) bool {
	if _, ok := s.m[v]; ok {
		return false
	}
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[v] = struct{}{}
	return true
}

// Has reports whether v is in the set. O(1) expected.
func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Remove deletes v and reports whether it was present. O(1) expected.
func (s *Set[T]) Remove(v T) bool {
	if _, ok := s.m[v]; !ok {
		return false
	}
	delete(s.m, v)
	return true
}

// All yields every value in an unspecified order.
func (s *Set[T]) All() iter.Seq[T] { return maps.Keys(s.m) }
