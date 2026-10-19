package text

import "unicode"

// IsPalindrome reports whether s reads the same forwards and backwards,
// ignoring case and every rune that is not a letter or digit. "A man, a
// plan, a canal: Panama" is a palindrome. The empty string is one.
//
// Time O(n). Space O(n) for the rune slice.
func IsPalindrome(s string) bool {
	var rs []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			rs = append(rs, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		if rs[i] != rs[j] {
			return false
		}
	}
	return true
}

// LongestPalindrome returns the longest substring of s that is an exact
// palindrome, comparing runes as they are. Ties go to the leftmost. It
// uses Manacher's algorithm.
//
// Time O(n). Space O(n).
func LongestPalindrome(s string) string {
	rs := []rune(s)
	if len(rs) < 2 {
		return s
	}
	// radius[i] is the palindrome radius centred at position i of the
	// virtual string #r0#r1#...#, where odd positions are runes and even
	// positions are gaps. It equals the palindrome's length in runes.
	m := 2*len(rs) + 1
	radius := make([]int, m)
	at := func(i int) rune {
		if i%2 == 0 {
			return -1
		}
		return rs[i/2]
	}
	center, right := 0, 0
	bestLen, bestCenter := 0, 0
	for i := range m {
		if i < right {
			radius[i] = min(right-i, radius[2*center-i])
		}
		for i-radius[i]-1 >= 0 && i+radius[i]+1 < m && at(i-radius[i]-1) == at(i+radius[i]+1) {
			radius[i]++
		}
		if i+radius[i] > right {
			center, right = i, i+radius[i]
		}
		if radius[i] > bestLen {
			bestLen, bestCenter = radius[i], i
		}
	}
	start := (bestCenter - bestLen) / 2
	return string(rs[start : start+bestLen])
}
