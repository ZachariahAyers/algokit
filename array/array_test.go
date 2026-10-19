package array

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestFlatten(t *testing.T) {
	got := Flatten([][]int{{1, 2}, nil, {}, {3}, {4, 5}})
	if !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("got %v", got)
	}
	if got := Flatten[int](nil); len(got) != 0 {
		t.Fatalf("nil: got %v", got)
	}
}

func TestSumAndProduct(t *testing.T) {
	if Sum([]int{}) != 0 || Sum([]int{1, -2, 3}) != 2 || Sum([]float64{0.5, 0.25}) != 0.75 {
		t.Fatal("Sum")
	}
	if Product([]int{}) != 1 || Product([]int{2, -3, 4}) != -24 || Product([]int{5, 0, 9}) != 0 {
		t.Fatal("Product")
	}
}

func TestMeanAndMedian(t *testing.T) {
	if _, err := Mean([]int{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Mean empty: %v", err)
	}
	if _, err := Median([]int(nil)); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Median nil: %v", err)
	}
	if m, _ := Mean([]int{1, 2, 3, 4}); m != 2.5 {
		t.Fatalf("Mean = %v", m)
	}
	// The float64 accumulator must not overflow where an int64 sum would.
	if m, _ := Mean([]int64{math.MaxInt64, math.MaxInt64}); m != float64(math.MaxInt64) {
		t.Fatalf("Mean overflow = %v", m)
	}
	in := []int{9, 1, 5}
	if m, _ := Median(in); m != 5 {
		t.Fatalf("Median odd = %v", m)
	}
	if !slices.Equal(in, []int{9, 1, 5}) {
		t.Fatal("Median modified its input")
	}
	if m, _ := Median([]int{4, 1, 3, 2}); m != 2.5 {
		t.Fatalf("Median even = %v", m)
	}
}

func TestDuplicatesAndIsUnique(t *testing.T) {
	got := Duplicates([]string{"a", "b", "a", "c", "b", "a"})
	if !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("got %q", got)
	}
	if Duplicates([]int{1, 2, 3}) != nil {
		t.Fatal("expected no duplicates")
	}
	if !IsUnique([]int{}) || !IsUnique([]int{1, 2, 3}) || IsUnique([]int{1, 2, 1}) {
		t.Fatal("IsUnique")
	}
}

func TestRotate(t *testing.T) {
	cases := []struct {
		k    int
		want []int
	}{
		{0, []int{1, 2, 3, 4, 5}},
		{2, []int{3, 4, 5, 1, 2}},
		{5, []int{1, 2, 3, 4, 5}},
		{7, []int{3, 4, 5, 1, 2}},
		{-1, []int{5, 1, 2, 3, 4}},
		{-6, []int{5, 1, 2, 3, 4}},
	}
	for _, c := range cases {
		s := []int{1, 2, 3, 4, 5}
		Rotate(s, c.k)
		if !slices.Equal(s, c.want) {
			t.Fatalf("Rotate(%d) = %v, want %v", c.k, s, c.want)
		}
	}
	Rotate([]int{}, 3)
}

func TestShuffleIsAPermutationAndUniform(t *testing.T) {
	r := rand.New(rand.NewPCG(9, 10))
	s := []int{0, 1, 2, 3, 4, 5, 6, 7}
	Shuffle(s, r)
	sorted := slices.Sorted(slices.Values(s))
	if !slices.Equal(sorted, []int{0, 1, 2, 3, 4, 5, 6, 7}) {
		t.Fatalf("not a permutation: %v", s)
	}
	// Each of the 6 orderings of 3 items should appear about 1/6 of the time.
	const trials = 60_000
	counts := map[[3]int]int{}
	for range trials {
		p := []int{0, 1, 2}
		Shuffle(p, r)
		counts[[3]int(p)]++
	}
	if len(counts) != 6 {
		t.Fatalf("saw %d orderings, want 6", len(counts))
	}
	for p, c := range counts {
		if c < trials/6*9/10 || c > trials/6*11/10 {
			t.Fatalf("ordering %v appeared %d times of %d", p, c, trials)
		}
	}
}

func TestMoveZeroes(t *testing.T) {
	s := []int{0, 1, 0, 3, 12, 0}
	MoveZeroes(s)
	if !slices.Equal(s, []int{1, 3, 12, 0, 0, 0}) {
		t.Fatalf("got %v", s)
	}
	strs := []string{"", "a", "", "b"}
	MoveZeroes(strs)
	if !slices.Equal(strs, []string{"a", "b", "", ""}) {
		t.Fatalf("got %q", strs)
	}
	MoveZeroes([]int(nil))
}
