package tree_test

import (
	"fmt"

	"github.com/ZachariahAyers/algokit/v2/tree"
)

func ExampleAVL() {
	var ages tree.AVL[string, int]
	ages.Put("Grace", 85)
	ages.Put("Ada", 36)
	ages.Put("Linus", 56)

	for name, age := range ages.InOrder() {
		fmt.Println(name, age)
	}
	// Output:
	// Ada 36
	// Grace 85
	// Linus 56
}
