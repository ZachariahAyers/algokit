// Package hashtable implements a hash map with separate chaining, to show
// how Go's built-in map works underneath. Use the built-in map in real
// code.
//
// Keys are hashed with [maphash.Comparable] under a per-table random seed,
// which stops an attacker who chooses keys from forcing every key into one
// bucket.
package hashtable

import (
	"hash/maphash"
	"iter"
)

const (
	minBuckets    = 8
	maxLoadFactor = 0.75
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

// Table maps keys to values. Build one with [New].
type Table[K comparable, V any] struct {
	buckets [][]entry[K, V]
	n       int
	seed    maphash.Seed
}

// New returns an empty table.
func New[K comparable, V any]() *Table[K, V] {
	return &Table[K, V]{buckets: make([][]entry[K, V], minBuckets), seed: maphash.MakeSeed()}
}

func (t *Table[K, V]) bucket(key K) int {
	return int(maphash.Comparable(t.seed, key) % uint64(len(t.buckets)))
}

// Len returns the number of keys. O(1).
func (t *Table[K, V]) Len() int { return t.n }

// Get returns the value stored under key. O(1) expected.
func (t *Table[K, V]) Get(key K) (V, bool) {
	for _, e := range t.buckets[t.bucket(key)] {
		if e.key == key {
			return e.value, true
		}
	}
	var zero V
	return zero, false
}

// Put stores value under key, replacing any existing value. It doubles the
// bucket count once the load factor passes 0.75. O(1) amortized expected.
func (t *Table[K, V]) Put(key K, value V) {
	b := t.bucket(key)
	for i := range t.buckets[b] {
		if t.buckets[b][i].key == key {
			t.buckets[b][i].value = value
			return
		}
	}
	t.buckets[b] = append(t.buckets[b], entry[K, V]{key, value})
	t.n++
	if float64(t.n) > maxLoadFactor*float64(len(t.buckets)) {
		t.resize(2 * len(t.buckets))
	}
}

// Delete removes key and reports whether it was present. O(1) expected.
func (t *Table[K, V]) Delete(key K) bool {
	b := t.bucket(key)
	chain := t.buckets[b]
	for i := range chain {
		if chain[i].key == key {
			last := len(chain) - 1
			chain[i] = chain[last]
			chain[last] = entry[K, V]{}
			t.buckets[b] = chain[:last]
			t.n--
			return true
		}
	}
	return false
}

// All yields every key and value in an unspecified order.
func (t *Table[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, chain := range t.buckets {
			for _, e := range chain {
				if !yield(e.key, e.value) {
					return
				}
			}
		}
	}
}

func (t *Table[K, V]) resize(size int) {
	old := t.buckets
	t.buckets = make([][]entry[K, V], size)
	for _, chain := range old {
		for _, e := range chain {
			b := t.bucket(e.key)
			t.buckets[b] = append(t.buckets[b], e)
		}
	}
}
