package dp_test

import (
	"fmt"

	"github.com/ZachariahAyers/algokit/v2/dp"
)

func ExampleKnapsack() {
	weights := []int{10, 20, 30}
	values := []int{60, 100, 120}
	best, chosen, _ := dp.Knapsack(weights, values, 50)
	fmt.Println(best, chosen)
	// Output: 220 [1 2]
}
