package search

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

type sortedSearch struct {
	name string
	fn   func([]int, int) (int, bool)
}

var sortedSearches = []sortedSearch{
	{"Binary", Binary[int]},
	{"Jump", Jump[int]},
	{"Interpolation", Interpolation[int]},
}

func TestSortedSearchesMatchStdlib(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 8))
	inputs := map[string][]int{
		"nil":        nil,
		"single":     {4},
		"duplicates": {1, 2, 2, 2, 2, 3, 9, 9},
		"all equal":  {5, 5, 5, 5, 5},
		"negatives":  {-40, -3, 0, 7, 100},
		"extremes":   {math.MinInt, -1, 0, 1, math.MaxInt},
	}
	random := make([]int, 1000)
	for i := range random {
		random[i] = rng.IntN(3000) - 1500
	}
	slices.Sort(random)
	inputs["random"] = random

	for _, alg := range sortedSearches {
		for name, s := range inputs {
			t.Run(alg.name+"/"+name, func(t *testing.T) {
				targets := []int{math.MinInt, math.MaxInt, -1, 0, 2, 3, 4, 6, 9, 10}
				targets = append(targets, s...)
				for _, target := range targets {
					wantI, wantOK := slices.BinarySearch(s, target)
					gotI, gotOK := alg.fn(s, target)
					if gotI != wantI || gotOK != wantOK {
						t.Fatalf("target %d: got (%d, %v), want (%d, %v)", target, gotI, gotOK, wantI, wantOK)
					}
				}
			})
		}
	}
}

func TestLinear(t *testing.T) {
	s := []string{"b", "a", "c", "a"}
	if i, ok := Linear(s, "a"); i != 1 || !ok {
		t.Fatalf("got (%d, %v), want (1, true)", i, ok)
	}
	if i, ok := Linear(s, "z"); i != len(s) || ok {
		t.Fatalf("got (%d, %v), want (%d, false)", i, ok, len(s))
	}
	if i, ok := Linear[int](nil, 1); i != 0 || ok {
		t.Fatalf("nil: got (%d, %v)", i, ok)
	}
}

func TestBinaryOnStrings(t *testing.T) {
	s := []string{"apple", "fig", "pear"}
	if i, ok := Binary(s, "grape"); i != 2 || ok {
		t.Fatalf("got (%d, %v), want (2, false)", i, ok)
	}
}
