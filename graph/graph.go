// Package graph implements a weighted graph, directed or undirected, on
// adjacency lists, with breadth-first and depth-first traversal.
//
// Vertices and each vertex's edges keep insertion order, so iteration and
// traversal order are deterministic.
package graph

import (
	"iter"
	"slices"
)

type edge[V comparable] struct {
	to     V
	weight float64
}

// Graph is a weighted graph over vertices of type V. Build one with [New].
type Graph[V comparable] struct {
	directed bool
	adj      map[V][]edge[V]
	order    []V
	edges    int
}

// New returns an empty graph. A directed graph stores each edge one way;
// an undirected graph stores it both ways.
func New[V comparable](directed bool) *Graph[V] {
	return &Graph[V]{directed: directed, adj: make(map[V][]edge[V])}
}

// Directed reports whether edges have a direction.
func (g *Graph[V]) Directed() bool { return g.directed }

// Order returns the number of vertices. O(1).
func (g *Graph[V]) Order() int { return len(g.order) }

// Size returns the number of edges. An undirected edge counts once. O(1).
func (g *Graph[V]) Size() int { return g.edges }

// AddVertex adds v and reports whether it was new. O(1).
func (g *Graph[V]) AddVertex(v V) bool {
	if _, ok := g.adj[v]; ok {
		return false
	}
	g.adj[v] = nil
	g.order = append(g.order, v)
	return true
}

// HasVertex reports whether v is in the graph. O(1).
func (g *Graph[V]) HasVertex(v V) bool {
	_, ok := g.adj[v]
	return ok
}

// RemoveVertex deletes v and every edge touching it, and reports whether v
// was present. O(V + E).
func (g *Graph[V]) RemoveVertex(v V) bool {
	out, ok := g.adj[v]
	if !ok {
		return false
	}
	// Every edge leaving v is listed once in out; in a directed graph the
	// edges entering v are counted as they are dropped below.
	g.edges -= len(out)
	delete(g.adj, v)
	g.order = slices.DeleteFunc(g.order, func(u V) bool { return u == v })
	for u, es := range g.adj {
		kept := slices.DeleteFunc(es, func(e edge[V]) bool { return e.to == v })
		if g.directed {
			g.edges -= len(es) - len(kept)
		}
		g.adj[u] = kept
	}
	return true
}

func (g *Graph[V]) find(u, v V) int {
	return slices.IndexFunc(g.adj[u], func(e edge[V]) bool { return e.to == v })
}

// AddEdge adds an edge from u to v with the given weight, adding either
// vertex if missing. If the edge exists its weight is replaced. O(deg u).
func (g *Graph[V]) AddEdge(u, v V, weight float64) {
	g.AddVertex(u)
	g.AddVertex(v)
	if i := g.find(u, v); i >= 0 {
		g.adj[u][i].weight = weight
		if !g.directed {
			g.adj[v][g.find(v, u)].weight = weight
		}
		return
	}
	g.adj[u] = append(g.adj[u], edge[V]{v, weight})
	if !g.directed && u != v {
		g.adj[v] = append(g.adj[v], edge[V]{u, weight})
	}
	g.edges++
}

// HasEdge reports whether an edge runs from u to v. O(deg u).
func (g *Graph[V]) HasEdge(u, v V) bool { return g.find(u, v) >= 0 }

// Weight returns the weight of the edge from u to v. O(deg u).
func (g *Graph[V]) Weight(u, v V) (float64, bool) {
	if i := g.find(u, v); i >= 0 {
		return g.adj[u][i].weight, true
	}
	return 0, false
}

// RemoveEdge deletes the edge from u to v and reports whether it existed.
// O(deg u + deg v).
func (g *Graph[V]) RemoveEdge(u, v V) bool {
	i := g.find(u, v)
	if i < 0 {
		return false
	}
	g.adj[u] = slices.Delete(g.adj[u], i, i+1)
	if !g.directed && u != v {
		j := g.find(v, u)
		g.adj[v] = slices.Delete(g.adj[v], j, j+1)
	}
	g.edges--
	return true
}

// Vertices yields every vertex in insertion order.
func (g *Graph[V]) Vertices() iter.Seq[V] { return slices.Values(g.order) }

// Neighbors yields each vertex reachable from v by one edge, with the
// edge weight, in edge insertion order.
func (g *Graph[V]) Neighbors(v V) iter.Seq2[V, float64] {
	return func(yield func(V, float64) bool) {
		for _, e := range g.adj[v] {
			if !yield(e.to, e.weight) {
				return
			}
		}
	}
}

// BFS yields the vertices reachable from start in breadth-first order,
// start first. It yields nothing if start is not in the graph.
//
// Time O(V + E). Space O(V).
func (g *Graph[V]) BFS(start V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if !g.HasVertex(start) {
			return
		}
		seen := map[V]bool{start: true}
		queue := []V{start}
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			if !yield(v) {
				return
			}
			for _, e := range g.adj[v] {
				if !seen[e.to] {
					seen[e.to] = true
					queue = append(queue, e.to)
				}
			}
		}
	}
}

// DFS yields the vertices reachable from start in depth-first preorder,
// visiting neighbours in edge insertion order, the same order a recursive
// search would produce. It uses an explicit stack, so deep graphs cannot
// overflow the call stack. It yields nothing if start is not in the graph.
//
// Time O(V + E). Space O(V + E).
func (g *Graph[V]) DFS(start V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if !g.HasVertex(start) {
			return
		}
		seen := map[V]bool{}
		stack := []V{start}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if seen[v] {
				continue
			}
			seen[v] = true
			if !yield(v) {
				return
			}
			es := g.adj[v]
			for i := len(es) - 1; i >= 0; i-- {
				if !seen[es[i].to] {
					stack = append(stack, es[i].to)
				}
			}
		}
	}
}
