// Package trie implements a prefix tree over the runes of strings. The
// zero value is an empty trie.
package trie

import (
	"iter"
	"maps"
	"slices"
)

type node struct {
	children map[rune]*node
	terminal bool // a stored word ends here
}

// Trie stores a set of strings and answers prefix queries.
type Trie struct {
	root node
	n    int
}

// Len returns the number of stored words. O(1).
func (t *Trie) Len() int { return t.n }

// Insert adds word and reports whether it was new. The empty string is a
// valid word. O(m) for m runes.
func (t *Trie) Insert(word string) bool {
	nd := &t.root
	for _, r := range word {
		child := nd.children[r]
		if child == nil {
			if nd.children == nil {
				nd.children = make(map[rune]*node)
			}
			child = &node{}
			nd.children[r] = child
		}
		nd = child
	}
	if nd.terminal {
		return false
	}
	nd.terminal = true
	t.n++
	return true
}

// find returns the node reached by walking prefix, or nil.
func (t *Trie) find(prefix string) *node {
	nd := &t.root
	for _, r := range prefix {
		if nd = nd.children[r]; nd == nil {
			return nil
		}
	}
	return nd
}

// Has reports whether word was inserted. O(m).
func (t *Trie) Has(word string) bool {
	nd := t.find(word)
	return nd != nil && nd.terminal
}

// HasPrefix reports whether any stored word starts with prefix. O(m).
func (t *Trie) HasPrefix(prefix string) bool {
	nd := t.find(prefix)
	return nd != nil && (nd.terminal || len(nd.children) > 0)
}

// Delete removes word and reports whether it was present, pruning nodes
// that no longer lead to a word. O(m).
func (t *Trie) Delete(word string) bool {
	rs := []rune(word)
	path := make([]*node, 0, len(rs)+1)
	nd := &t.root
	path = append(path, nd)
	for _, r := range rs {
		if nd = nd.children[r]; nd == nil {
			return false
		}
		path = append(path, nd)
	}
	if !nd.terminal {
		return false
	}
	nd.terminal = false
	t.n--
	for i := len(rs); i > 0; i-- {
		child := path[i]
		if child.terminal || len(child.children) > 0 {
			break
		}
		delete(path[i-1].children, rs[i-1])
	}
	return true
}

// WithPrefix yields every stored word that starts with prefix, in
// ascending rune order. O(m + output size · k log k) for k children
// per node, sorted on the way down.
func (t *Trie) WithPrefix(prefix string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if nd := t.find(prefix); nd != nil {
			walk(nd, []rune(prefix), yield)
		}
	}
}

// All yields every stored word in ascending rune order.
func (t *Trie) All() iter.Seq[string] { return t.WithPrefix("") }

func walk(nd *node, buf []rune, yield func(string) bool) bool {
	if nd.terminal && !yield(string(buf)) {
		return false
	}
	for _, r := range slices.Sorted(maps.Keys(nd.children)) {
		if !walk(nd.children[r], append(buf, r), yield) {
			return false
		}
	}
	return true
}
