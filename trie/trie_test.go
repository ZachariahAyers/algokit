package trie

import (
	"slices"
	"testing"
)

func TestInsertHasDelete(t *testing.T) {
	var tr Trie
	for _, w := range []string{"car", "card", "care", "cat", "dog", ""} {
		if !tr.Insert(w) {
			t.Fatalf("Insert(%q) reported duplicate", w)
		}
	}
	if tr.Insert("car") {
		t.Fatal("second Insert(car) reported new")
	}
	if tr.Len() != 6 {
		t.Fatalf("Len = %d", tr.Len())
	}
	if !tr.Has("car") || tr.Has("ca") || !tr.Has("") || tr.Has("cards") {
		t.Fatal("Has")
	}
	if !tr.HasPrefix("ca") || !tr.HasPrefix("card") || tr.HasPrefix("cb") || !tr.HasPrefix("") {
		t.Fatal("HasPrefix")
	}
	if !tr.Delete("car") || tr.Delete("car") || tr.Delete("ca") || tr.Delete("zebra") {
		t.Fatal("Delete results")
	}
	if tr.Has("car") || !tr.Has("card") || tr.Len() != 5 {
		t.Fatal("Delete removed the wrong words")
	}
	if !tr.Delete("dog") || tr.HasPrefix("d") {
		t.Fatal("Delete should prune the dog branch")
	}
	if len(tr.root.children) != 1 {
		t.Fatalf("root children = %d after pruning, want 1", len(tr.root.children))
	}
}

func TestWithPrefixIsSorted(t *testing.T) {
	var tr Trie
	for _, w := range []string{"tea", "ten", "to", "inn", "tėst", "te", "i"} {
		tr.Insert(w)
	}
	if got := slices.Collect(tr.WithPrefix("te")); !slices.Equal(got, []string{"te", "tea", "ten"}) {
		t.Fatalf("WithPrefix(te) = %q", got)
	}
	if got := slices.Collect(tr.All()); !slices.Equal(got, []string{"i", "inn", "te", "tea", "ten", "to", "tėst"}) {
		t.Fatalf("All = %q", got)
	}
	if got := slices.Collect(tr.WithPrefix("x")); got != nil {
		t.Fatalf("missing prefix = %q", got)
	}
	n := 0
	for range tr.All() {
		n++
		break
	}
	if n != 1 {
		t.Fatal("early exit")
	}
}
