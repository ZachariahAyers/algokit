package combinatorics

import (
	"slices"
	"testing"
)

func equal2D[T comparable](a, b [][]T) bool {
	return slices.EqualFunc(a, b, slices.Equal[[]T])
}

func TestCartesianProduct(t *testing.T) {
	got := slices.Collect(CartesianProduct([]int{1, 2}, []int{3, 4, 5}))
	want := [][]int{{1, 3}, {1, 4}, {1, 5}, {2, 3}, {2, 4}, {2, 5}}
	if !equal2D(got, want) {
		t.Fatalf("got %v", got)
	}
	if got := slices.Collect(CartesianProduct[int]()); !equal2D(got, [][]int{{}}) {
		t.Fatalf("no sets: got %v", got)
	}
	if got := slices.Collect(CartesianProduct([]int{1}, nil)); got != nil {
		t.Fatalf("empty set: got %v", got)
	}
}

func TestCombinations(t *testing.T) {
	got := slices.Collect(Combinations([]string{"a", "b", "c", "d"}, 2))
	want := [][]string{{"a", "b"}, {"a", "c"}, {"a", "d"}, {"b", "c"}, {"b", "d"}, {"c", "d"}}
	if !equal2D(got, want) {
		t.Fatalf("got %v", got)
	}
	if got := slices.Collect(Combinations([]int{1, 2}, 0)); !equal2D(got, [][]int{{}}) {
		t.Fatalf("k=0: got %v", got)
	}
	for _, k := range []int{-1, 3} {
		if got := slices.Collect(Combinations([]int{1, 2}, k)); got != nil {
			t.Fatalf("k=%d: got %v", k, got)
		}
	}
	// C(10, 4) = 210.
	if n := len(slices.Collect(Combinations(make([]int, 10), 4))); n != 210 {
		t.Fatalf("C(10,4) = %d", n)
	}
}

func TestPowerSet(t *testing.T) {
	got := slices.Collect(PowerSet([]int{1, 2, 3}))
	want := [][]int{{}, {1}, {2}, {3}, {1, 2}, {1, 3}, {2, 3}, {1, 2, 3}}
	if !equal2D(got, want) {
		t.Fatalf("got %v", got)
	}
	if got := slices.Collect(PowerSet[int](nil)); !equal2D(got, [][]int{{}}) {
		t.Fatalf("empty: got %v", got)
	}
}

func TestPermutations(t *testing.T) {
	got := slices.Collect(Permutations([]int{3, 1, 2}))
	want := [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}
	if !equal2D(got, want) {
		t.Fatalf("got %v", got)
	}
	got = slices.Collect(Permutations([]int{1, 1, 2}))
	want = [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}}
	if !equal2D(got, want) {
		t.Fatalf("multiset: got %v", got)
	}
	if got := slices.Collect(Permutations[int](nil)); !equal2D(got, [][]int{{}}) {
		t.Fatalf("empty: got %v", got)
	}
	in := []int{2, 1}
	_ = slices.Collect(Permutations(in))
	if !slices.Equal(in, []int{2, 1}) {
		t.Fatal("input was modified")
	}
}

func TestIteratorsStopEarly(t *testing.T) {
	n := 0
	for range PowerSet(make([]int, 40)) {
		n++
		if n == 5 {
			break
		}
	}
	if n != 5 {
		t.Fatalf("got %d", n)
	}
	for range Permutations([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}) {
		break
	}
	for range CartesianProduct([]int{1, 2}, []int{3}) {
		break
	}
}

func TestCombinationSum(t *testing.T) {
	got := CombinationSum([]int{2, 3, 6, 7}, 7)
	if !equal2D(got, [][]int{{2, 2, 3}, {7}}) {
		t.Fatalf("got %v", got)
	}
	got = CombinationSum([]int{3, 2, 2, 5, -1, 0}, 8)
	if !equal2D(got, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}) {
		t.Fatalf("got %v", got)
	}
	if got := CombinationSum([]int{2}, 1); got != nil {
		t.Fatalf("unreachable: got %v", got)
	}
	if got := CombinationSum([]int{1}, 0); got != nil {
		t.Fatalf("target 0: got %v", got)
	}
}
