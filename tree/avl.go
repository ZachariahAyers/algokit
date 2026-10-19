package tree

import "cmp"

// AVL is a self-balancing binary search tree: the heights of every node's
// two subtrees differ by at most one, so the height stays below
// 1.44·log2(n+2). Get, Put and Delete are O(log n) in the worst case.
type AVL[K cmp.Ordered, V any] struct {
	core[K, V]
}

// Put stores value under key, replacing any existing value. O(log n).
func (t *AVL[K, V]) Put(key K, value V) {
	var added bool
	t.root, added = avlPut(t.root, key, value)
	if added {
		t.n++
	}
}

// Delete removes key and reports whether it was present. O(log n).
func (t *AVL[K, V]) Delete(key K) bool {
	var removed bool
	t.root, removed = avlDelete(t.root, key)
	if removed {
		t.n--
	}
	return removed
}

func avlPut[K cmp.Ordered, V any](nd *node[K, V], key K, value V) (*node[K, V], bool) {
	if nd == nil {
		return &node[K, V]{key: key, value: value, height: 1}, true
	}
	var added bool
	switch c := cmp.Compare(key, nd.key); {
	case c < 0:
		nd.left, added = avlPut(nd.left, key, value)
	case c > 0:
		nd.right, added = avlPut(nd.right, key, value)
	default:
		nd.value = value
		return nd, false
	}
	return rebalance(nd), added
}

func avlDelete[K cmp.Ordered, V any](nd *node[K, V], key K) (*node[K, V], bool) {
	if nd == nil {
		return nil, false
	}
	var removed bool
	switch c := cmp.Compare(key, nd.key); {
	case c < 0:
		nd.left, removed = avlDelete(nd.left, key)
	case c > 0:
		nd.right, removed = avlDelete(nd.right, key)
	default:
		if nd.left == nil {
			return nd.right, true
		}
		if nd.right == nil {
			return nd.left, true
		}
		succ := minNode(nd.right)
		nd.key, nd.value = succ.key, succ.value
		nd.right, _ = avlDelete(nd.right, succ.key)
		removed = true
	}
	return rebalance(nd), removed
}

func height[K cmp.Ordered, V any](nd *node[K, V]) int {
	if nd == nil {
		return 0
	}
	return nd.height
}

func fix[K cmp.Ordered, V any](nd *node[K, V]) {
	nd.height = 1 + max(height(nd.left), height(nd.right))
}

func balance[K cmp.Ordered, V any](nd *node[K, V]) int {
	return height(nd.left) - height(nd.right)
}

// rebalance updates nd's height and applies the single or double rotation
// that restores the AVL property, returning the new subtree root.
func rebalance[K cmp.Ordered, V any](nd *node[K, V]) *node[K, V] {
	fix(nd)
	switch b := balance(nd); {
	case b > 1:
		if balance(nd.left) < 0 {
			nd.left = rotateLeft(nd.left)
		}
		return rotateRight(nd)
	case b < -1:
		if balance(nd.right) > 0 {
			nd.right = rotateRight(nd.right)
		}
		return rotateLeft(nd)
	}
	return nd
}

func rotateRight[K cmp.Ordered, V any](nd *node[K, V]) *node[K, V] {
	l := nd.left
	nd.left, l.right = l.right, nd
	fix(nd)
	fix(l)
	return l
}

func rotateLeft[K cmp.Ordered, V any](nd *node[K, V]) *node[K, V] {
	r := nd.right
	nd.right, r.left = r.left, nd
	fix(nd)
	fix(r)
	return r
}
