package linkedlist

import "iter"

type doublyNode[T any] struct {
	value      T
	prev, next *doublyNode[T]
}

// Doubly is a doubly linked list. Unlike [Singly], PopBack is O(1) and
// index operations walk from whichever end is nearer.
type Doubly[T any] struct {
	head, tail *doublyNode[T]
	n          int
}

// Len returns the number of elements. O(1).
func (l *Doubly[T]) Len() int { return l.n }

// Front returns the first element. O(1).
func (l *Doubly[T]) Front() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// Back returns the last element. O(1).
func (l *Doubly[T]) Back() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// PushFront inserts v at the front. O(1).
func (l *Doubly[T]) PushFront(v T) { l.insertBefore(l.head, v) }

// PushBack appends v. O(1).
func (l *Doubly[T]) PushBack(v T) { l.insertBefore(nil, v) }

// PopFront removes and returns the first element. O(1).
func (l *Doubly[T]) PopFront() (T, bool) { return l.unlink(l.head) }

// PopBack removes and returns the last element. O(1).
func (l *Doubly[T]) PopBack() (T, bool) { return l.unlink(l.tail) }

// insertBefore links a new node holding v before at, or at the back when
// at is nil.
func (l *Doubly[T]) insertBefore(at *doublyNode[T], v T) {
	node := &doublyNode[T]{value: v, next: at}
	if at == nil {
		node.prev = l.tail
		l.tail = node
	} else {
		node.prev = at.prev
		at.prev = node
	}
	if node.prev == nil {
		l.head = node
	} else {
		node.prev.next = node
	}
	l.n++
}

// unlink removes node and returns its value; a nil node reports false.
func (l *Doubly[T]) unlink(node *doublyNode[T]) (T, bool) {
	if node == nil {
		var zero T
		return zero, false
	}
	if node.prev == nil {
		l.head = node.next
	} else {
		node.prev.next = node.next
	}
	if node.next == nil {
		l.tail = node.prev
	} else {
		node.next.prev = node.prev
	}
	l.n--
	return node.value, true
}

// nodeAt returns the node at index i, or nil outside [0, n). O(min(i, n-i)).
func (l *Doubly[T]) nodeAt(i int) *doublyNode[T] {
	if i < 0 || i >= l.n {
		return nil
	}
	if i < l.n/2 {
		node := l.head
		for range i {
			node = node.next
		}
		return node
	}
	node := l.tail
	for range l.n - 1 - i {
		node = node.prev
	}
	return node
}

// At returns the element at index i. O(min(i, n-i)).
func (l *Doubly[T]) At(i int) (T, bool) {
	node := l.nodeAt(i)
	if node == nil {
		var zero T
		return zero, false
	}
	return node.value, true
}

// Set replaces the element at index i. O(min(i, n-i)).
func (l *Doubly[T]) Set(i int, v T) bool {
	node := l.nodeAt(i)
	if node == nil {
		return false
	}
	node.value = v
	return true
}

// InsertAt inserts v so that it ends up at index i, for i in [0, Len()].
// O(min(i, n-i)).
func (l *Doubly[T]) InsertAt(i int, v T) bool {
	if i < 0 || i > l.n {
		return false
	}
	l.insertBefore(l.nodeAt(i), v)
	return true
}

// RemoveAt removes and returns the element at index i. O(min(i, n-i)).
func (l *Doubly[T]) RemoveAt(i int) (T, bool) { return l.unlink(l.nodeAt(i)) }

// Reverse reverses the list in place by swapping each node's links. O(n).
func (l *Doubly[T]) Reverse() {
	for node := l.head; node != nil; node = node.prev {
		node.prev, node.next = node.next, node.prev
	}
	l.head, l.tail = l.tail, l.head
}

// Clear removes every element. O(1).
func (l *Doubly[T]) Clear() { *l = Doubly[T]{} }

// All yields the index and value of each element, front to back.
func (l *Doubly[T]) All() iter.Seq2[int, T] {
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

// Backward yields the index and value of each element, back to front.
func (l *Doubly[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := l.n - 1
		for node := l.tail; node != nil; node = node.prev {
			if !yield(i, node.value) {
				return
			}
			i--
		}
	}
}
