package dp

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

func isSubsequence[T comparable](sub, s []T) bool {
	i := 0
	for _, v := range s {
		if i < len(sub) && sub[i] == v {
			i++
		}
	}
	return i == len(sub)
}

func TestLongestCommonSubsequence(t *testing.T) {
	a, b := []rune("ABCBDAB"), []rune("BDCABA")
	got := LongestCommonSubsequence(a, b)
	if len(got) != 4 || !isSubsequence(got, a) || !isSubsequence(got, b) {
		t.Fatalf("got %q", string(got))
	}
	if got := LongestCommonSubsequence([]int{1, 2}, []int{3}); len(got) != 0 {
		t.Fatalf("disjoint: got %v", got)
	}
	if got := LongestCommonSubsequence[int](nil, []int{1}); len(got) != 0 {
		t.Fatalf("nil: got %v", got)
	}
}

func TestShortestCommonSupersequence(t *testing.T) {
	rng := rand.New(rand.NewPCG(13, 14))
	for range 300 {
		a := randomRunes(rng, 8)
		b := randomRunes(rng, 8)
		got := ShortestCommonSupersequence(a, b)
		lcs := len(LongestCommonSubsequence(a, b))
		if len(got) != len(a)+len(b)-lcs || !isSubsequence(a, got) || !isSubsequence(b, got) {
			t.Fatalf("SCS(%q, %q) = %q", string(a), string(b), string(got))
		}
	}
	if got := string(ShortestCommonSupersequence([]rune("abac"), []rune("cab"))); len(got) != 5 {
		t.Fatalf("got %q", got)
	}
}

func randomRunes(rng *rand.Rand, maxLen int) []rune {
	rs := make([]rune, rng.IntN(maxLen+1))
	for i := range rs {
		rs[i] = rune('a' + rng.IntN(3))
	}
	return rs
}

// bruteLIS returns the length of the longest strictly increasing
// subsequence by checking every subset.
func bruteLIS(s []int, less func(a, b int) bool) int {
	best := 0
	for mask := range 1 << len(s) {
		var picked []int
		for i := range s {
			if mask&(1<<i) != 0 {
				picked = append(picked, s[i])
			}
		}
		ok := true
		for i := 1; i < len(picked); i++ {
			if !less(picked[i-1], picked[i]) {
				ok = false
				break
			}
		}
		if ok {
			best = max(best, len(picked))
		}
	}
	return best
}

func TestLongestMonotoneSubsequences(t *testing.T) {
	got := LongestIncreasingSubsequence([]int{10, 9, 2, 5, 3, 7, 101, 18})
	if !slices.Equal(got, []int{2, 3, 7, 18}) {
		t.Fatalf("LIS = %v", got)
	}
	if got := LongestIncreasingSubsequence([]int{}); len(got) != 0 {
		t.Fatalf("empty: %v", got)
	}
	if got := LongestIncreasingSubsequence([]int{4, 4, 4}); !slices.Equal(got, []int{4}) {
		t.Fatalf("equal values must not count as increasing: %v", got)
	}
	if got := LongestDecreasingSubsequence([]int{5, 1, 4, 2, 3, 0}); len(got) != 4 {
		t.Fatalf("LDS = %v", got)
	}

	inc := func(a, b int) bool { return a < b }
	dec := func(a, b int) bool { return a > b }
	rng := rand.New(rand.NewPCG(15, 16))
	for range 300 {
		s := make([]int, rng.IntN(13))
		for i := range s {
			s[i] = rng.IntN(6)
		}
		lis := LongestIncreasingSubsequence(s)
		if len(lis) != bruteLIS(s, inc) || !isSubsequence(lis, s) || !slices.IsSortedFunc(lis, func(a, b int) int { return a - b }) {
			t.Fatalf("LIS(%v) = %v", s, lis)
		}
		lds := LongestDecreasingSubsequence(s)
		if len(lds) != bruteLIS(s, dec) || !isSubsequence(lds, s) {
			t.Fatalf("LDS(%v) = %v", s, lds)
		}
		for i := 1; i < len(lds); i++ {
			if lds[i] >= lds[i-1] {
				t.Fatalf("LDS(%v) = %v is not strictly decreasing", s, lds)
			}
		}
	}
}

func TestEditDistance(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"intention", "execution", 5},
		{"café", "cafe", 1},
	}
	for _, c := range cases {
		if got := EditDistance([]rune(c.a), []rune(c.b)); got != c.want {
			t.Errorf("EditDistance(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
		if got := EditDistance([]rune(c.b), []rune(c.a)); got != c.want {
			t.Errorf("EditDistance is not symmetric for %q, %q", c.a, c.b)
		}
	}
}

func TestKnapsack(t *testing.T) {
	best, chosen, err := Knapsack([]int{1, 3, 4, 5}, []int{1, 4, 5, 7}, 7)
	if err != nil || best != 9 || !slices.Equal(chosen, []int{1, 2}) {
		t.Fatalf("got %d %v %v", best, chosen, err)
	}
	best, chosen, err = Knapsack([]int{10, 20, 30}, []int{60, 100, 120}, 50)
	if err != nil || best != 220 || !slices.Equal(chosen, []int{1, 2}) {
		t.Fatalf("got %d %v %v", best, chosen, err)
	}
	if best, chosen, err := Knapsack(nil, nil, 10); err != nil || best != 0 || chosen != nil {
		t.Fatalf("empty: %d %v %v", best, chosen, err)
	}
	if best, _, _ := Knapsack([]int{5}, []int{10}, 4); best != 0 {
		t.Fatalf("too heavy: %d", best)
	}
	if best, chosen, _ := Knapsack([]int{0, 2}, []int{3, 4}, 0); best != 3 || !slices.Equal(chosen, []int{0}) {
		t.Fatalf("zero weight: %d %v", best, chosen)
	}
	for _, c := range []struct {
		w, v []int
		c    int
	}{{[]int{1}, nil, 1}, {[]int{-1}, []int{1}, 1}, {nil, nil, -1}} {
		if _, _, err := Knapsack(c.w, c.v, c.c); !errors.Is(err, ErrKnapsackInput) {
			t.Fatalf("Knapsack(%v, %v, %d) error = %v", c.w, c.v, c.c, err)
		}
	}
}
