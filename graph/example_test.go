package graph_test

import (
	"fmt"
	"slices"

	"github.com/ZachariahAyers/algokit/v2/graph"
)

func ExampleGraph_BFS() {
	g := graph.New[string](false)
	g.AddEdge("home", "shop", 2)
	g.AddEdge("home", "park", 1)
	g.AddEdge("park", "lake", 3)

	fmt.Println(slices.Collect(g.BFS("home")))
	// Output: [home shop park lake]
}
