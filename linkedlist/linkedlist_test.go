package linkedlist

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// list is the API both list types share, so one model test covers both.
type list interface {
	Len() int
	Front() (int, bool)
	Back() (int, bool)
	PushFront(int)
	PushBack(int)
	PopFront() (int, bool)
	PopBack() (int, bool)
	At(int) (int, bool)
	Set(int, int) bool
	InsertAt(int, int) bool
	RemoveAt(int) (int, bool)
	Reverse()
	Clear()
}

func values(l interface{ Len() int }, at func(int) (int, bool)) []int {
	out := make([]int, 0, l.Len())
	for i := range l.Len() {
		v, _ := at(i)
		out = append(out, v)
	}
	return out
}

func TestListsMatchSliceModel(t *testing.T) {
	impls := map[string]func() list{
		"Singly": func() list { return &Singly[int]{} },
		"Doubly": func() list { return &Doubly[int]{} },
	}
	for name, mk := range impls {
		t.Run(name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(17, 18))
			l := mk()
			var model []int
			for step := range 5000 {
				v := rng.IntN(1000)
				i := rng.IntN(len(model)+3) - 1 // includes out-of-range -1 and len+1
				switch op := rng.IntN(10); op {
				case 0:
					l.PushFront(v)
					model = slices.Insert(model, 0, v)
				case 1:
					l.PushBack(v)
					model = append(model, v)
				case 2:
					got, ok := l.PopFront()
					if ok != (len(model) > 0) || ok && got != model[0] {
						t.Fatalf("step %d PopFront = %d, %v", step, got, ok)
					}
					if ok {
						model = model[1:]
					}
				case 3:
					got, ok := l.PopBack()
					if ok != (len(model) > 0) || ok && got != model[len(model)-1] {
						t.Fatalf("step %d PopBack = %d, %v", step, got, ok)
					}
					if ok {
						model = model[:len(model)-1]
					}
				case 4:
					ok := l.InsertAt(i, v)
					valid := i >= 0 && i <= len(model)
					if ok != valid {
						t.Fatalf("step %d InsertAt(%d) = %v", step, i, ok)
					}
					if valid {
						model = slices.Insert(model, i, v)
					}
				case 5:
					got, ok := l.RemoveAt(i)
					valid := i >= 0 && i < len(model)
					if ok != valid || valid && got != model[i] {
						t.Fatalf("step %d RemoveAt(%d) = %d, %v", step, i, got, ok)
					}
					if valid {
						model = slices.Delete(model, i, i+1)
					}
				case 6:
					ok := l.Set(i, v)
					valid := i >= 0 && i < len(model)
					if ok != valid {
						t.Fatalf("step %d Set(%d) = %v", step, i, ok)
					}
					if valid {
						model[i] = v
					}
				case 7:
					l.Reverse()
					slices.Reverse(model)
				case 8:
					if rng.IntN(50) == 0 {
						l.Clear()
						model = nil
					}
				default:
					got, ok := l.At(i)
					valid := i >= 0 && i < len(model)
					if ok != valid || valid && got != model[i] {
						t.Fatalf("step %d At(%d) = %d, %v", step, i, got, ok)
					}
				}
				if l.Len() != len(model) {
					t.Fatalf("step %d Len = %d, want %d", step, l.Len(), len(model))
				}
				front, okF := l.Front()
				back, okB := l.Back()
				if len(model) == 0 {
					if okF || okB {
						t.Fatalf("step %d empty list reports an end", step)
					}
				} else if front != model[0] || back != model[len(model)-1] {
					t.Fatalf("step %d ends = %d, %d; want %d, %d", step, front, back, model[0], model[len(model)-1])
				}
			}
			if got := values(l, l.At); !slices.Equal(got, model) {
				t.Fatalf("final contents %v, want %v", got, model)
			}
		})
	}
}

func TestIterators(t *testing.T) {
	var s Singly[string]
	var d Doubly[string]
	for _, v := range []string{"a", "b", "c"} {
		s.PushBack(v)
		d.PushBack(v)
	}
	var fromS, fromD, back []string
	for i, v := range s.All() {
		if i != len(fromS) {
			t.Fatalf("Singly index %d", i)
		}
		fromS = append(fromS, v)
	}
	for _, v := range d.All() {
		fromD = append(fromD, v)
	}
	for i, v := range d.Backward() {
		if i != 2-len(back) {
			t.Fatalf("Backward index %d", i)
		}
		back = append(back, v)
	}
	if !slices.Equal(fromS, []string{"a", "b", "c"}) || !slices.Equal(fromD, fromS) || !slices.Equal(back, []string{"c", "b", "a"}) {
		t.Fatalf("got %v %v %v", fromS, fromD, back)
	}
	for range s.All() {
		break
	}
	for range d.Backward() {
		break
	}
}
