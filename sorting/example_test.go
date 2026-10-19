package sorting_test

import (
	"fmt"

	"github.com/ZachariahAyers/algokit/v2/sorting"
)

func ExampleMerge() {
	s := []int{5, 2, 9, 1, 5}
	sorting.Merge(s)
	fmt.Println(s)
	// Output: [1 2 5 5 9]
}
