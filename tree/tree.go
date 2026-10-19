// Package tree implements ordered maps on binary search trees: a plain
// unbalanced [BST] and a self-balancing [AVL] tree.
//
// Both share the read side (lookups, min and max, and the four
// traversals); they differ only in how Put and Delete reshape the tree.
// The zero value of each is an empty tree.
package tree

import (
	"cmp"
	"iter"
)

type node[K cmp.Ordered, V any] struct {
	key         K
	value       V
	left, right *node[K, V]
	height      int // maintained by AVL only; a leaf has height 1
}

// core holds the read operations both trees share.
type core[K cmp.Ordered, V any] struct {
	root *node[K, V]
	n    int
}

// Len returns the number of keys. O(1).
func (t *core[K, V]) Len() int { return t.n }

// Get returns the value stored under key. O(h) for tree height h.
func (t *core[K, V]) Get(key K) (V, bool) {
	for nd := t.root; nd != nil; {
		switch c := cmp.Compare(key, nd.key); {
		case c < 0:
			nd = nd.left
		case c > 0:
			nd = nd.right
		default:
			return nd.value, true
		}
	}
	var zero V
	return zero, false
}

// Has reports whether key is present. O(h).
func (t *core[K, V]) Has(key K) bool {
	_, ok := t.Get(key)
	return ok
}

// Min returns the smallest key and its value. O(h).
func (t *core[K, V]) Min() (K, V, bool) {
	if t.root == nil {
		var k K
		var v V
		return k, v, false
	}
	nd := minNode(t.root)
	return nd.key, nd.value, true
}

// Max returns the largest key and its value. O(h).
func (t *core[K, V]) Max() (K, V, bool) {
	nd := t.root
	if nd == nil {
		var k K
		var v V
		return k, v, false
	}
	for nd.right != nil {
		nd = nd.right
	}
	return nd.key, nd.value, true
}

// Clear removes every key. O(1).
func (t *core[K, V]) Clear() { *t = core[K, V]{} }

// InOrder yields keys in ascending order. O(n) total, O(h) stack.
func (t *core[K, V]) InOrder() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) { inOrder(t.root, yield) }
}

// PreOrder yields each node before its subtrees: depth first, root first.
func (t *core[K, V]) PreOrder() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) { preOrder(t.root, yield) }
}

// PostOrder yields each node after its subtrees.
func (t *core[K, V]) PostOrder() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) { postOrder(t.root, yield) }
}

// LevelOrder yields nodes breadth first, left to right within each level.
// O(n) total, O(w) queue for the widest level w.
func (t *core[K, V]) LevelOrder() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if t.root == nil {
			return
		}
		level := []*node[K, V]{t.root}
		for len(level) > 0 {
			var next []*node[K, V]
			for _, nd := range level {
				if !yield(nd.key, nd.value) {
					return
				}
				if nd.left != nil {
					next = append(next, nd.left)
				}
				if nd.right != nil {
					next = append(next, nd.right)
				}
			}
			level = next
		}
	}
}

// Each traversal returns false once yield asks to stop.
func inOrder[K cmp.Ordered, V any](nd *node[K, V], yield func(K, V) bool) bool {
	return nd == nil || inOrder(nd.left, yield) && yield(nd.key, nd.value) && inOrder(nd.right, yield)
}

func preOrder[K cmp.Ordered, V any](nd *node[K, V], yield func(K, V) bool) bool {
	return nd == nil || yield(nd.key, nd.value) && preOrder(nd.left, yield) && preOrder(nd.right, yield)
}

func postOrder[K cmp.Ordered, V any](nd *node[K, V], yield func(K, V) bool) bool {
	return nd == nil || postOrder(nd.left, yield) && postOrder(nd.right, yield) && yield(nd.key, nd.value)
}

func minNode[K cmp.Ordered, V any](nd *node[K, V]) *node[K, V] {
	for nd.left != nil {
		nd = nd.left
	}
	return nd
}
