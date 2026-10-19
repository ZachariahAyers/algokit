package text

import (
	"slices"
	"strings"
	"unicode"

	"github.com/ZachariahAyers/algokit/v2/combinatorics"
)

// Anagrams returns every distinct rearrangement of the runes of s, in
// ascending order, including s itself. A word with repeated letters yields
// no duplicates: "aab" gives "aab", "aba", "baa".
//
// Time O(n·n!) in the worst case. Space O(n·n!) for the result.
func Anagrams(s string) []string {
	var out []string
	for p := range combinatorics.Permutations([]rune(s)) {
		out = append(out, string(p))
	}
	return out
}

// CharFrequency counts how often each rune occurs in s.
//
// Time O(n). Space O(k) for k distinct runes.
func CharFrequency(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, r := range s {
		freq[r]++
	}
	return freq
}

// Reverse returns s with its runes in reverse order.
//
// Time O(n). Space O(n).
func Reverse(s string) string {
	rs := []rune(s)
	slices.Reverse(rs)
	return string(rs)
}

// ReverseWords returns the whitespace-separated words of s in reverse
// order, joined by single spaces.
//
// Time O(n). Space O(n).
func ReverseWords(s string) string {
	words := strings.Fields(s)
	slices.Reverse(words)
	return strings.Join(words, " ")
}

// IsRotation reports whether b is a rotation of a, such as "erbottlewat"
// of "waterbottle", by checking that b occurs in a+a.
//
// Time O(n) expected for the substring search. Space O(n).
func IsRotation(a, b string) bool {
	return len(a) == len(b) && strings.Contains(a+a, b)
}

// RemoveVowels returns s without the ASCII vowels a, e, i, o and u in
// either case.
//
// Time O(n). Space O(n).
func RemoveVowels(s string) string {
	return strings.Map(func(r rune) rune {
		switch unicode.ToLower(r) {
		case 'a', 'e', 'i', 'o', 'u':
			return -1
		}
		return r
	}, s)
}
