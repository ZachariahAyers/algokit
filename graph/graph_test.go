package graph

import (
	"slices"
	"testing"
)

func TestUndirectedEdges(t *testing.T) {
	g := New[string](false)
	g.AddEdge("a", "b", 1)
	g.AddEdge("b", "c", 2)
	g.AddEdge("a", "a", 5)
	if g.Order() != 3 || g.Size() != 3 || g.Directed() {
		t.Fatalf("Order %d Size %d", g.Order(), g.Size())
	}
	if !g.HasEdge("b", "a") || !g.HasEdge("a", "a") || g.HasEdge("a", "c") {
		t.Fatal("HasEdge")
	}
	g.AddEdge("b", "a", 9)
	if w, _ := g.Weight("a", "b"); w != 9 || g.Size() != 3 {
		t.Fatalf("replacing a weight: w=%v size=%d", w, g.Size())
	}
	if !g.RemoveEdge("b", "a") || g.HasEdge("a", "b") || g.RemoveEdge("a", "b") || g.Size() != 2 {
		t.Fatal("RemoveEdge")
	}
	if _, ok := g.Weight("a", "b"); ok {
		t.Fatal("Weight of removed edge")
	}
	if !g.RemoveVertex("b") || g.RemoveVertex("b") || g.HasVertex("b") || g.Size() != 1 || g.Order() != 2 {
		t.Fatalf("RemoveVertex: order %d size %d", g.Order(), g.Size())
	}
	if !g.RemoveVertex("a") || g.Size() != 0 {
		t.Fatalf("self-loop removal left size %d", g.Size())
	}
}

func TestDirectedEdges(t *testing.T) {
	g := New[int](true)
	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 1, 1)
	g.AddEdge(3, 1, 1)
	g.AddEdge(1, 1, 1)
	if g.Size() != 4 || !g.HasEdge(3, 1) || g.HasEdge(1, 3) {
		t.Fatal("directed edges")
	}
	if g.AddVertex(1) || !g.AddVertex(4) {
		t.Fatal("AddVertex")
	}
	if !g.RemoveVertex(1) || g.Size() != 0 || g.Order() != 3 {
		t.Fatalf("RemoveVertex: order %d size %d", g.Order(), g.Size())
	}
	if got := slices.Collect(g.Vertices()); !slices.Equal(got, []int{2, 3, 4}) {
		t.Fatalf("Vertices = %v", got)
	}
}

// sample builds
//
//	1 - 2 - 4
//	|   |
//	3 - 5   6 - 7
func sample() *Graph[int] {
	g := New[int](false)
	for _, e := range [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 5}, {6, 7}} {
		g.AddEdge(e[0], e[1], 1)
	}
	return g
}

func TestTraversals(t *testing.T) {
	g := sample()
	if got := slices.Collect(g.BFS(1)); !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("BFS = %v", got)
	}
	if got := slices.Collect(g.DFS(1)); !slices.Equal(got, []int{1, 2, 4, 5, 3}) {
		t.Fatalf("DFS = %v", got)
	}
	if got := slices.Collect(g.BFS(6)); !slices.Equal(got, []int{6, 7}) {
		t.Fatalf("BFS from 6 = %v", got)
	}
	if got := slices.Collect(g.DFS(99)); got != nil {
		t.Fatalf("DFS from missing vertex = %v", got)
	}
	var neighbors []int
	for v, w := range g.Neighbors(2) {
		if w != 1 {
			t.Fatalf("weight %v", w)
		}
		neighbors = append(neighbors, v)
	}
	if !slices.Equal(neighbors, []int{1, 4, 5}) {
		t.Fatalf("Neighbors = %v", neighbors)
	}
	for _, seq := range [](func(func(int) bool)){g.BFS(1), g.DFS(1)} {
		n := 0
		for range seq {
			n++
			break
		}
		if n != 1 {
			t.Fatal("early exit")
		}
	}
}

func TestDFSHandlesDeepPaths(t *testing.T) {
	g := New[int](true)
	const n = 200_000
	for i := range n - 1 {
		g.AddEdge(i, i+1, 1)
	}
	count := 0
	for range g.DFS(0) {
		count++
	}
	if count != n {
		t.Fatalf("visited %d of %d", count, n)
	}
}
