package stack

import (
	"slices"
	"testing"
)

func TestStack(t *testing.T) {
	var s Stack[int]
	if !s.IsEmpty() {
		t.Fatal("zero value should be empty")
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty")
	}
	if _, ok := s.Peek(); ok {
		t.Fatal("Peek on empty")
	}
	for i := range 5 {
		s.Push(i)
	}
	if top, _ := s.Peek(); top != 4 || s.Len() != 5 {
		t.Fatalf("Peek = %d, Len = %d", top, s.Len())
	}
	if got := slices.Collect(s.All()); !slices.Equal(got, []int{4, 3, 2, 1, 0}) {
		t.Fatalf("All = %v", got)
	}
	for want := 4; want >= 0; want-- {
		if got, ok := s.Pop(); !ok || got != want {
			t.Fatalf("Pop = %d, %v; want %d", got, ok, want)
		}
	}
	s.Push(9)
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear")
	}
}
