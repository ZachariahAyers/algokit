package tree

import (
	"cmp"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

type orderedMap interface {
	Len() int
	Get(int) (string, bool)
	Has(int) bool
	Put(int, string)
	Delete(int) bool
}

// checkOrder verifies the search-tree property and returns the node count.
func checkOrder[V any](t *testing.T, nd *node[int, V], lo, hi int) int {
	t.Helper()
	if nd == nil {
		return 0
	}
	if nd.key <= lo || nd.key >= hi {
		t.Fatalf("key %d outside (%d, %d)", nd.key, lo, hi)
	}
	return 1 + checkOrder(t, nd.left, lo, nd.key) + checkOrder(t, nd.right, nd.key, hi)
}

// checkAVL verifies stored heights and the balance bound, returning height.
func checkAVL[V any](t *testing.T, nd *node[int, V]) int {
	t.Helper()
	if nd == nil {
		return 0
	}
	l, r := checkAVL(t, nd.left), checkAVL(t, nd.right)
	if l-r > 1 || r-l > 1 {
		t.Fatalf("node %d unbalanced: %d vs %d", nd.key, l, r)
	}
	if h := 1 + max(l, r); nd.height != h {
		t.Fatalf("node %d stores height %d, actual %d", nd.key, nd.height, h)
	}
	return nd.height
}

func TestTreesMatchMapModel(t *testing.T) {
	bst, avl := &BST[int, string]{}, &AVL[int, string]{}
	trees := map[string]struct {
		m    orderedMap
		root func() *node[int, string]
	}{
		"BST": {bst, func() *node[int, string] { return bst.root }},
		"AVL": {avl, func() *node[int, string] { return avl.root }},
	}
	for name, tr := range trees {
		t.Run(name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(23, 24))
			model := map[int]string{}
			for step := range 4000 {
				k := rng.IntN(300)
				switch rng.IntN(3) {
				case 0, 1:
					v := string(rune('a' + rng.IntN(26)))
					tr.m.Put(k, v)
					model[k] = v
				default:
					_, want := model[k]
					if got := tr.m.Delete(k); got != want {
						t.Fatalf("step %d Delete(%d) = %v", step, k, got)
					}
					delete(model, k)
				}
				got, ok := tr.m.Get(k)
				want, wantOK := model[k]
				if ok != wantOK || got != want || tr.m.Has(k) != wantOK {
					t.Fatalf("step %d Get(%d) = %q, %v", step, k, got, ok)
				}
				if tr.m.Len() != len(model) {
					t.Fatalf("step %d Len = %d, want %d", step, tr.m.Len(), len(model))
				}
				if step%100 == 0 {
					if n := checkOrder(t, tr.root(), math.MinInt, math.MaxInt); n != len(model) {
						t.Fatalf("step %d tree holds %d nodes, want %d", step, n, len(model))
					}
					if name == "AVL" {
						checkAVL(t, tr.root())
					}
				}
			}
		})
	}
}

func TestAVLStaysShallowOnSortedInput(t *testing.T) {
	var avl AVL[int, struct{}]
	const n = 100_000
	for i := range n {
		avl.Put(i, struct{}{})
	}
	h := checkAVL(t, avl.root)
	if limit := 1.44 * math.Log2(n+2); float64(h) > limit {
		t.Fatalf("height %d exceeds %.1f", h, limit)
	}
}

func TestTraversals(t *testing.T) {
	//        4
	//      /   \
	//     2     6
	//    / \   / \
	//   1   3 5   7
	var tr BST[int, int]
	for _, k := range []int{4, 2, 6, 1, 3, 5, 7} {
		tr.Put(k, k*10)
	}
	keys := func(seq func(func(int, int) bool)) []int {
		var out []int
		for k := range seq {
			out = append(out, k)
		}
		return out
	}
	cases := map[string]struct {
		got, want []int
	}{
		"InOrder":    {keys(tr.InOrder()), []int{1, 2, 3, 4, 5, 6, 7}},
		"PreOrder":   {keys(tr.PreOrder()), []int{4, 2, 1, 3, 6, 5, 7}},
		"PostOrder":  {keys(tr.PostOrder()), []int{1, 3, 2, 5, 7, 6, 4}},
		"LevelOrder": {keys(tr.LevelOrder()), []int{4, 2, 6, 1, 3, 5, 7}},
	}
	for name, c := range cases {
		if !slices.Equal(c.got, c.want) {
			t.Errorf("%s = %v, want %v", name, c.got, c.want)
		}
	}
	if k, v, ok := tr.Min(); k != 1 || v != 10 || !ok {
		t.Fatalf("Min = %d %d %v", k, v, ok)
	}
	if k, _, ok := tr.Max(); k != 7 || !ok {
		t.Fatalf("Max = %d %v", k, ok)
	}
	// Early exit from every traversal.
	for _, seq := range []func(func(int, int) bool){tr.InOrder(), tr.PreOrder(), tr.PostOrder(), tr.LevelOrder()} {
		n := 0
		for range seq {
			n++
			if n == 2 {
				break
			}
		}
		if n != 2 {
			t.Fatalf("early exit visited %d", n)
		}
	}
	tr.Clear()
	if _, _, ok := tr.Min(); ok || tr.Len() != 0 {
		t.Fatal("Clear")
	}
	if _, _, ok := tr.Max(); ok {
		t.Fatal("Max on empty")
	}
}

func TestInOrderIsSorted(t *testing.T) {
	var avl AVL[string, int]
	words := []string{"pear", "apple", "fig", "kiwi", "banana", "apple"}
	for i, w := range words {
		avl.Put(w, i)
	}
	got := map[string]int{}
	var order []string
	for k, v := range avl.InOrder() {
		order = append(order, k)
		got[k] = v
	}
	if !slices.IsSortedFunc(order, cmp.Compare[string]) || len(order) != 5 {
		t.Fatalf("order = %v", order)
	}
	if !maps.Equal(got, map[string]int{"apple": 5, "banana": 4, "fig": 2, "kiwi": 3, "pear": 0}) {
		t.Fatalf("values = %v", got)
	}
}
