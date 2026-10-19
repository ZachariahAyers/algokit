package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

type intSort struct {
	name string
	fn   func([]int)
}

var comparisonSorts = []intSort{
	{"Bubble", Bubble[int]},
	{"Insertion", Insertion[int]},
	{"Selection", Selection[int]},
	{"Shell", Shell[int]},
	{"Merge", Merge[int]},
	{"Quick", Quick[int]},
	{"Heap", Heap[int]},
	{"Counting", Counting[int]},
}

func edgeCases() map[string][]int {
	rng := rand.New(rand.NewPCG(1, 2))
	random := make([]int, 2000)
	for i := range random {
		random[i] = rng.IntN(4000) - 2000
	}
	fewKeys := make([]int, 2000)
	for i := range fewKeys {
		fewKeys[i] = rng.IntN(3)
	}
	ascending := make([]int, 500)
	descending := make([]int, 500)
	for i := range ascending {
		ascending[i] = i
		descending[i] = 500 - i
	}
	return map[string][]int{
		"nil":        nil,
		"empty":      {},
		"single":     {7},
		"pair":       {2, 1},
		"equal":      {5, 5, 5, 5},
		"negatives":  {3, -1, 0, -7, 2},
		"ascending":  ascending,
		"descending": descending,
		"few keys":   fewKeys,
		"random":     random,
	}
}

func TestSortsAgreeWithStdlib(t *testing.T) {
	for _, alg := range comparisonSorts {
		for name, in := range edgeCases() {
			t.Run(alg.name+"/"+name, func(t *testing.T) {
				got := slices.Clone(in)
				want := slices.Clone(in)
				slices.Sort(want)
				alg.fn(got)
				if !slices.Equal(got, want) {
					t.Fatalf("got %v, want %v", got, want)
				}
			})
		}
	}
}

func TestSortsHandleStrings(t *testing.T) {
	in := []string{"pear", "apple", "fig", "", "apple", "Banana"}
	want := slices.Clone(in)
	slices.Sort(want)
	for _, fn := range []func([]string){Bubble[string], Insertion[string], Selection[string], Shell[string], Merge[string], Quick[string], Heap[string]} {
		got := slices.Clone(in)
		fn(got)
		if !slices.Equal(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
}

// Equal cmp.Ordered values are indistinguishable except for floating-point
// zeros: -0 == +0 under <, yet the sign bit tells them apart. A stable sort
// must keep the zeros' sign sequence exactly as it was.
func TestStableSortsKeepEqualKeysInOrder(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	in := make([]float64, 400)
	var signs []bool
	for i := range in {
		switch rng.IntN(3) {
		case 0:
			in[i] = math.Copysign(0, -1)
			signs = append(signs, true)
		case 1:
			in[i] = 0
			signs = append(signs, false)
		default:
			in[i] = float64((rng.IntN(10) + 1) * (1 - 2*rng.IntN(2))) // nonzero, either sign
		}
	}
	stable := map[string]func([]float64){
		"Bubble":    Bubble[float64],
		"Insertion": Insertion[float64],
		"Merge":     Merge[float64],
	}
	for name, fn := range stable {
		t.Run(name, func(t *testing.T) {
			got := slices.Clone(in)
			fn(got)
			if !slices.IsSorted(got) {
				t.Fatalf("not sorted: %v", got)
			}
			var gotSigns []bool
			for _, v := range got {
				if v == 0 {
					gotSigns = append(gotSigns, math.Signbit(v))
				}
			}
			if !slices.Equal(gotSigns, signs) {
				t.Fatal("equal keys changed relative order")
			}
		})
	}
}

func TestRadixAgreesWithStdlib(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 6))
	cases := map[string][]uint64{
		"nil":    nil,
		"single": {9},
		"zeros":  {0, 0, 0},
		"max":    {^uint64(0), 0, 1 << 63, 42},
	}
	random := make([]uint64, 3000)
	for i := range random {
		random[i] = rng.Uint64() >> rng.UintN(64)
	}
	cases["random"] = random
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			got := slices.Clone(in)
			want := slices.Clone(in)
			slices.Sort(want)
			Radix(got)
			if !slices.Equal(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
	bytes := []uint8{200, 3, 255, 0, 3}
	Radix(bytes)
	if !slices.Equal(bytes, []uint8{0, 3, 3, 200, 255}) {
		t.Fatalf("uint8: got %v", bytes)
	}
}

func TestCountingHandlesExtremeNarrowRange(t *testing.T) {
	s := []int8{127, -128, 0, 127, -128}
	Counting(s)
	if !slices.Equal(s, []int8{-128, -128, 0, 127, 127}) {
		t.Fatalf("got %v", s)
	}
}

func TestCountingPanicsOnHugeRange(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a range wider than MaxCountingRange")
		}
	}()
	Counting([]int64{math.MinInt64, math.MaxInt64})
}
