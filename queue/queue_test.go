package queue

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestQueueFIFOAcrossWraparound(t *testing.T) {
	var q Queue[int]
	if _, ok := q.Pop(); ok {
		t.Fatal("Pop on empty")
	}
	if _, ok := q.Front(); ok {
		t.Fatal("Front on empty")
	}
	if _, ok := q.Back(); ok {
		t.Fatal("Back on empty")
	}
	rng := rand.New(rand.NewPCG(19, 20))
	var model []int
	next := 0
	for range 10_000 {
		if rng.IntN(3) > 0 {
			q.Push(next)
			model = append(model, next)
			next++
		} else {
			got, ok := q.Pop()
			if ok != (len(model) > 0) || ok && got != model[0] {
				t.Fatalf("Pop = %d, %v", got, ok)
			}
			if ok {
				model = model[1:]
			}
		}
		if q.Len() != len(model) || q.IsEmpty() != (len(model) == 0) {
			t.Fatalf("Len = %d, want %d", q.Len(), len(model))
		}
		if len(model) > 0 {
			f, _ := q.Front()
			b, _ := q.Back()
			if f != model[0] || b != model[len(model)-1] {
				t.Fatalf("Front/Back = %d/%d", f, b)
			}
		}
	}
	if got := slices.Collect(q.All()); !slices.Equal(got, model) {
		t.Fatalf("All = %v, want %v", got, model)
	}
	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("Clear")
	}
}
