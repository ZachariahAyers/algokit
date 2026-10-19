// Package linkedlist implements generic singly and doubly linked lists.
//
// Both lists keep head and tail pointers and a length, so Len, Front,
// Back, PushFront and PushBack are O(1). Index operations walk the list.
// Index arguments outside the valid range make the call report false
// instead of panicking. The zero value of each list is an empty list.
package linkedlist

import "iter"

type singlyNode[T any] struct {
	value T
	next  *singlyNode[T]
}

// Singly is a singly linked list.
type Singly[T any] struct {
	head, tail *singlyNode[T]
	n          int
}

// Len returns the number of elements. O(1).
func (l *Singly[T]) Len() int { return l.n }

// Front returns the first element. O(1).
func (l *Singly[T]) Front() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// Back returns the last element. O(1).
func (l *Singly[T]) Back() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// PushFront inserts v at the front. O(1).
func (l *Singly[T]) PushFront(v T) {
	l.head = &singlyNode[T]{value: v, next: l.head}
	if l.tail == nil {
		l.tail = l.head
	}
	l.n++
}

// PushBack appends v. O(1).
func (l *Singly[T]) PushBack(v T) {
	node := &singlyNode[T]{value: v}
	if l.tail == nil {
		l.head = node
	} else {
		l.tail.next = node
	}
	l.tail = node
	l.n++
}

// PopFront removes and returns the first element. O(1).
func (l *Singly[T]) PopFront() (T, bool) {
	return l.RemoveAt(0)
}

// PopBack removes and returns the last element. O(n): a singly linked
// list must walk to the node before the tail.
func (l *Singly[T]) PopBack() (T, bool) {
	return l.RemoveAt(l.n - 1)
}

// nodeAt returns the node at index i, which must be in [0, n). O(i).
func (l *Singly[T]) nodeAt(i int) *singlyNode[T] {
	node := l.head
	for range i {
		node = node.next
	}
	return node
}

// At returns the element at index i. O(i).
func (l *Singly[T]) At(i int) (T, bool) {
	if i < 0 || i >= l.n {
		var zero T
		return zero, false
	}
	return l.nodeAt(i).value, true
}

// Set replaces the element at index i. O(i).
func (l *Singly[T]) Set(i int, v T) bool {
	if i < 0 || i >= l.n {
		return false
	}
	l.nodeAt(i).value = v
	return true
}

// InsertAt inserts v so that it ends up at index i, for i in [0, Len()].
// O(i).
func (l *Singly[T]) InsertAt(i int, v T) bool {
	switch {
	case i < 0 || i > l.n:
		return false
	case i == 0:
		l.PushFront(v)
	case i == l.n:
		l.PushBack(v)
	default:
		prev := l.nodeAt(i - 1)
		prev.next = &singlyNode[T]{value: v, next: prev.next}
		l.n++
	}
	return true
}

// RemoveAt removes and returns the element at index i. O(i).
func (l *Singly[T]) RemoveAt(i int) (T, bool) {
	if i < 0 || i >= l.n {
		var zero T
		return zero, false
	}
	var removed *singlyNode[T]
	if i == 0 {
		removed = l.head
		l.head = removed.next
		if l.head == nil {
			l.tail = nil
		}
	} else {
		prev := l.nodeAt(i - 1)
		removed = prev.next
		prev.next = removed.next
		if removed == l.tail {
			l.tail = prev
		}
	}
	l.n--
	return removed.value, true
}

// Reverse reverses the list in place by relinking nodes. O(n).
func (l *Singly[T]) Reverse() {
	var prev *singlyNode[T]
	l.tail = l.head
	for node := l.head; node != nil; {
		next := node.next
		node.next = prev
		prev, node = node, next
	}
	l.head = prev
}

// Clear removes every element. O(1).
func (l *Singly[T]) Clear() { *l = Singly[T]{} }

// All yields the index and value of each element, front to back.
func (l *Singly[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for node := l.head; node != nil; node = node.next {
			if !yield(i, node.value) {
				return
			}
			i++
		}
	}
}
