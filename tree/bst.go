package tree

import "cmp"

// BST is an unbalanced binary search tree. Operations cost O(h) where h is
// the tree height: O(log n) for random insertion order, O(n) for sorted
// insertion order. Use [AVL] when the order is not known to be random.
type BST[K cmp.Ordered, V any] struct {
	core[K, V]
}

// Put stores value under key, replacing any existing value. O(h).
func (t *BST[K, V]) Put(key K, value V) {
	link := &t.root
	for *link != nil {
		switch c := cmp.Compare(key, (*link).key); {
		case c < 0:
			link = &(*link).left
		case c > 0:
			link = &(*link).right
		default:
			(*link).value = value
			return
		}
	}
	*link = &node[K, V]{key: key, value: value}
	t.n++
}

// Delete removes key and reports whether it was present. A node with two
// children is replaced by its in-order successor. O(h).
func (t *BST[K, V]) Delete(key K) bool {
	link := &t.root
	for *link != nil {
		nd := *link
		switch c := cmp.Compare(key, nd.key); {
		case c < 0:
			link = &nd.left
		case c > 0:
			link = &nd.right
		default:
			switch {
			case nd.left == nil:
				*link = nd.right
			case nd.right == nil:
				*link = nd.left
			default:
				// Unlink the successor from the right subtree and put it
				// in nd's place.
				succLink := &nd.right
				for (*succLink).left != nil {
					succLink = &(*succLink).left
				}
				succ := *succLink
				*succLink = succ.right
				succ.left, succ.right = nd.left, nd.right
				*link = succ
			}
			t.n--
			return true
		}
	}
	return false
}
