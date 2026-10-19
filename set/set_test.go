package set

import (
	"slices"
	"testing"
)

func TestSet(t *testing.T) {
	var s Set[string]
	if s.Has("a") || s.Remove("a") || s.Len() != 0 {
		t.Fatal("zero value should be empty")
	}
	if !s.Add("a") || s.Add("a") || !s.Add("b") {
		t.Fatal("Add results")
	}
	if !s.Has("a") || s.Len() != 2 {
		t.Fatal("Has or Len")
	}
	if !s.Remove("a") || s.Remove("a") || s.Has("a") || s.Len() != 1 {
		t.Fatal("Remove")
	}
	o := Of(3, 1, 3, 2)
	if got := slices.Sorted(o.All()); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Of/All = %v", got)
	}
}
