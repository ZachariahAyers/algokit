package numeric

import (
	"fmt"
	"strings"
)

var romanTable = []struct {
	value  int
	symbol string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

// ToRoman returns the canonical Roman numeral for n, greedily taking the
// largest symbol that fits.
//
// Time O(1): at most 15 symbols for n <= 3999. Returns [ErrDomain] outside
// 1..3999.
func ToRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		return "", fmt.Errorf("%w: %d is outside 1..3999", ErrDomain, n)
	}
	var b strings.Builder
	for _, r := range romanTable {
		for n >= r.value {
			b.WriteString(r.symbol)
			n -= r.value
		}
	}
	return b.String(), nil
}

// FromRoman parses a canonical, upper-case Roman numeral. Non-canonical
// forms such as "IIII" or "IC" are rejected, so FromRoman and [ToRoman]
// are exact inverses.
//
// Time O(len(s)). Returns [ErrDomain] for an empty or invalid numeral.
func FromRoman(s string) (int, error) {
	n, rest := 0, s
	for _, r := range romanTable {
		for strings.HasPrefix(rest, r.symbol) {
			n += r.value
			rest = rest[len(r.symbol):]
		}
	}
	if rest != "" || n == 0 {
		return 0, fmt.Errorf("%w: %q is not a Roman numeral", ErrDomain, s)
	}
	// The greedy parse accepts some non-canonical strings ("IIII", "VV");
	// a round trip rejects them.
	if canon, _ := ToRoman(n); canon != s {
		return 0, fmt.Errorf("%w: %q is not canonical (want %q)", ErrDomain, s, canon)
	}
	return n, nil
}
