package numeric

import "strings"

var (
	ones = []string{
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		"seventeen", "eighteen", "nineteen",
	}
	tens   = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	scales = []string{"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion"}
)

// ToWords spells n in English using the short scale, without "and":
// 1_234 is "one thousand two hundred thirty-four" and -5 is "minus five".
// It covers the full int64 range.
//
// Time O(d) for d digits. Space O(d).
func ToWords(n int64) string {
	if n == 0 {
		return ones[0]
	}
	// Work in uint64 so math.MinInt64 negates without overflow.
	u := uint64(n)
	negative := n < 0
	if negative {
		u = -u
	}
	var groups []string
	for scale := 0; u > 0; scale++ {
		if g := u % 1000; g > 0 {
			part := underThousand(int(g))
			if scales[scale] != "" {
				part += " " + scales[scale]
			}
			groups = append(groups, part)
		}
		u /= 1000
	}
	var b strings.Builder
	if negative {
		b.WriteString("minus ")
	}
	for i := len(groups) - 1; i >= 0; i-- {
		b.WriteString(groups[i])
		if i > 0 {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

func underThousand(n int) string {
	var parts []string
	if n >= 100 {
		parts = append(parts, ones[n/100]+" hundred")
		n %= 100
	}
	switch {
	case n >= 20 && n%10 != 0:
		parts = append(parts, tens[n/10]+"-"+ones[n%10])
	case n >= 20:
		parts = append(parts, tens[n/10])
	case n > 0:
		parts = append(parts, ones[n])
	}
	return strings.Join(parts, " ")
}
