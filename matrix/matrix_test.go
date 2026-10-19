package matrix

import (
	"errors"
	"slices"
	"testing"
)

func TestRotate(t *testing.T) {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	if err := Rotate(m); err != nil {
		t.Fatal(err)
	}
	want := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}
	if !slices.EqualFunc(m, want, slices.Equal) {
		t.Fatalf("got %v", m)
	}
}

func TestRotateFourTimesIsIdentity(t *testing.T) {
	m := [][]string{{"a", "b"}, {"c", "d"}}
	for range 4 {
		if err := Rotate(m); err != nil {
			t.Fatal(err)
		}
	}
	if !slices.EqualFunc(m, [][]string{{"a", "b"}, {"c", "d"}}, slices.Equal) {
		t.Fatalf("got %v", m)
	}
}

func TestRotateEdgeCases(t *testing.T) {
	if err := Rotate[int](nil); err != nil {
		t.Fatalf("empty: %v", err)
	}
	one := [][]int{{42}}
	if err := Rotate(one); err != nil || one[0][0] != 42 {
		t.Fatalf("1x1: %v %v", one, err)
	}
	ragged := [][]int{{1, 2}, {3}}
	if err := Rotate(ragged); !errors.Is(err, ErrNotSquare) {
		t.Fatalf("ragged: %v", err)
	}
	if !slices.Equal(ragged[0], []int{1, 2}) {
		t.Fatal("ragged matrix was modified")
	}
	if err := Rotate([][]int{{1, 2, 3}, {4, 5, 6}}); !errors.Is(err, ErrNotSquare) {
		t.Fatalf("2x3: %v", err)
	}
}
