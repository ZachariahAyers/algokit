// Package queue implements a generic first-in, first-out queue on a ring
// buffer that doubles when full. The zero value is an empty queue.
package queue

import "iter"

// Queue is a FIFO queue.
type Queue[T any] struct {
	buf  []T
	head int
	n    int
}

// Len returns the number of elements. O(1).
func (q *Queue[T]) Len() int { return q.n }

// IsEmpty reports whether the queue has no elements. O(1).
func (q *Queue[T]) IsEmpty() bool { return q.n == 0 }

// Push adds v at the back. O(1) amortized.
func (q *Queue[T]) Push(v T) {
	if q.n == len(q.buf) {
		q.grow()
	}
	q.buf[(q.head+q.n)%len(q.buf)] = v
	q.n++
}

// Pop removes and returns the front element. O(1).
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if q.n == 0 {
		return zero, false
	}
	v := q.buf[q.head]
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.n--
	return v, true
}

// Front returns the element Pop would return next. O(1).
func (q *Queue[T]) Front() (T, bool) {
	if q.n == 0 {
		var zero T
		return zero, false
	}
	return q.buf[q.head], true
}

// Back returns the most recently pushed element. O(1).
func (q *Queue[T]) Back() (T, bool) {
	if q.n == 0 {
		var zero T
		return zero, false
	}
	return q.buf[(q.head+q.n-1)%len(q.buf)], true
}

// Clear removes every element. O(1).
func (q *Queue[T]) Clear() { *q = Queue[T]{} }

// All yields the elements from front to back.
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range q.n {
			if !yield(q.buf[(q.head+i)%len(q.buf)]) {
				return
			}
		}
	}
}

// grow doubles the buffer and unwraps the elements to start at index 0.
func (q *Queue[T]) grow() {
	buf := make([]T, max(4, 2*len(q.buf)))
	for i := range q.n {
		buf[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf, q.head = buf, 0
}
