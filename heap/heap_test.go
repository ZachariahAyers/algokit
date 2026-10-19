package heap

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func drain[T any](h *Heap[T]) []T {
	var out []T
	for !h.IsEmpty() {
		v, _ := h.Pop()
		out = append(out, v)
	}
	return out
}

func TestMinAndMaxHeaps(t *testing.T) {
	rng := rand.New(rand.NewPCG(21, 22))
	in := make([]int, 1000)
	for i := range in {
		in[i] = rng.IntN(200) - 100
	}
	minH, maxH := NewMin[int](), NewMax[int]()
	for _, v := range in {
		minH.Push(v)
		maxH.Push(v)
	}
	asc := slices.Sorted(slices.Values(in))
	desc := slices.Clone(asc)
	slices.Reverse(desc)
	if top, _ := minH.Peek(); top != asc[0] {
		t.Fatalf("min Peek = %d", top)
	}
	if got := drain(minH); !slices.Equal(got, asc) {
		t.Fatal("min-heap pop order is not ascending")
	}
	if got := drain(maxH); !slices.Equal(got, desc) {
		t.Fatal("max-heap pop order is not descending")
	}
	if _, ok := minH.Pop(); ok {
		t.Fatal("Pop on empty")
	}
	if _, ok := minH.Peek(); ok {
		t.Fatal("Peek on empty")
	}
}

func TestFromHeapifiesAndPriorityQueue(t *testing.T) {
	type task struct {
		name     string
		priority int
	}
	tasks := []task{{"write", 2}, {"test", 1}, {"ship", 3}, {"plan", 0}}
	pq := From(tasks, func(a, b task) bool { return a.priority < b.priority })
	if pq.Len() != 4 {
		t.Fatalf("Len = %d", pq.Len())
	}
	var order []string
	for _, tk := range drain(pq) {
		order = append(order, tk.name)
	}
	if !slices.Equal(order, []string{"plan", "test", "write", "ship"}) {
		t.Fatalf("order = %v", order)
	}
	h := From([]int{5, 3, 8}, func(a, b int) bool { return a < b })
	if got := slices.Sorted(h.All()); !slices.Equal(got, []int{3, 5, 8}) {
		t.Fatalf("All = %v", got)
	}
}
