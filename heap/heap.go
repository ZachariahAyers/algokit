// Package heap implements a generic binary heap, which also serves as a
// priority queue. The ordering is a less function, so one type covers
// min-heaps, max-heaps and heaps keyed on any field. For an in-place
// heapsort see sorting.Heap.
package heap

import (
	"cmp"
	"iter"
)

// Heap is a binary heap whose root is the element that orders first under
// less. Build one with [New], [NewMin] or [NewMax].
type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

// New returns an empty heap ordered by less.
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

// NewMin returns an empty heap whose root is the smallest element.
func NewMin[T cmp.Ordered]() *Heap[T] {
	return New(func(a, b T) bool { return a < b })
}

// NewMax returns an empty heap whose root is the largest element.
func NewMax[T cmp.Ordered]() *Heap[T] {
	return New(func(a, b T) bool { return a > b })
}

// From returns a heap ordered by less holding items, built bottom-up.
// The heap takes ownership of items.
//
// Time O(n), which beats n pushes at O(n log n).
func From[T any](items []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{items: items, less: less}
	for i := len(items)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h
}

// Len returns the number of elements. O(1).
func (h *Heap[T]) Len() int { return len(h.items) }

// IsEmpty reports whether the heap has no elements. O(1).
func (h *Heap[T]) IsEmpty() bool { return len(h.items) == 0 }

// Push adds v. O(log n).
func (h *Heap[T]) Push(v T) {
	h.items = append(h.items, v)
	h.up(len(h.items) - 1)
}

// Peek returns the root without removing it. O(1).
func (h *Heap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

// Pop removes and returns the root. O(log n).
func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	n := len(h.items)
	if n == 0 {
		return zero, false
	}
	root := h.items[0]
	h.items[0] = h.items[n-1]
	h.items[n-1] = zero
	h.items = h.items[:n-1]
	if n > 1 {
		h.down(0)
	}
	return root, true
}

// All yields the elements in heap-array order, which is not sorted order.
// Pop repeatedly for sorted output.
func (h *Heap[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range h.items {
			if !yield(v) {
				return
			}
		}
	}
}

func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.items[i], h.items[parent]) {
			return
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

func (h *Heap[T]) down(i int) {
	n := len(h.items)
	for {
		child := 2*i + 1
		if child >= n {
			return
		}
		if r := child + 1; r < n && h.less(h.items[r], h.items[child]) {
			child = r
		}
		if !h.less(h.items[child], h.items[i]) {
			return
		}
		h.items[i], h.items[child] = h.items[child], h.items[i]
		i = child
	}
}
