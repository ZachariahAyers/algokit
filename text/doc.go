// Package text implements string algorithms: palindromes, anagrams,
// character statistics, word and case transforms, and simple encodings.
//
// Functions work on Unicode code points (runes), not bytes, unless their
// documentation says otherwise. They do not segment grapheme clusters, so a
// letter written with a combining accent counts as two runes.
package text
