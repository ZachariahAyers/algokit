package hashtable

import (
	"maps"
	"math/rand/v2"
	"testing"
)

func TestTableMatchesBuiltinMap(t *testing.T) {
	rng := rand.New(rand.NewPCG(25, 26))
	tbl := New[int, int]()
	model := map[int]int{}
	for step := range 20_000 {
		k := rng.IntN(2000)
		if rng.IntN(3) == 0 {
			_, want := model[k]
			if got := tbl.Delete(k); got != want {
				t.Fatalf("step %d Delete(%d) = %v", step, k, got)
			}
			delete(model, k)
		} else {
			tbl.Put(k, step)
			model[k] = step
		}
		got, ok := tbl.Get(k)
		want, wantOK := model[k]
		if ok != wantOK || got != want {
			t.Fatalf("step %d Get(%d) = %d, %v", step, k, got, ok)
		}
	}
	if tbl.Len() != len(model) {
		t.Fatalf("Len = %d, want %d", tbl.Len(), len(model))
	}
	if !maps.Equal(maps.Collect(tbl.All()), model) {
		t.Fatal("All does not match the model")
	}
	if float64(tbl.Len()) > maxLoadFactor*float64(len(tbl.buckets)) {
		t.Fatalf("load factor %.2f exceeds the limit", float64(tbl.Len())/float64(len(tbl.buckets)))
	}
}

func TestStructKeysAndEarlyExit(t *testing.T) {
	type point struct{ x, y int }
	tbl := New[point, string]()
	tbl.Put(point{1, 2}, "a")
	tbl.Put(point{1, 2}, "b")
	tbl.Put(point{2, 1}, "c")
	if v, _ := tbl.Get(point{1, 2}); v != "b" || tbl.Len() != 2 {
		t.Fatalf("Get = %q, Len = %d", v, tbl.Len())
	}
	if _, ok := tbl.Get(point{9, 9}); ok {
		t.Fatal("missing key found")
	}
	n := 0
	for range tbl.All() {
		n++
		break
	}
	if n != 1 {
		t.Fatal("early exit")
	}
}
