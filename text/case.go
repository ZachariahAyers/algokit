package text

import (
	"strings"
	"unicode"
)

// Words splits s into words for case conversion. A word boundary is any
// rune that is not a letter or digit, a lower-case letter or digit
// followed by an upper-case letter ("fooBar"), and the last capital of an
// acronym followed by a lower-case letter ("HTTPServer" is "HTTP",
// "Server").
//
// Time O(n). Space O(n).
func Words(s string) []string {
	rs := []rune(s)
	var words []string
	start := -1
	flush := func(end int) {
		if start >= 0 {
			words = append(words, string(rs[start:end]))
			start = -1
		}
	}
	for i, r := range rs {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			flush(i)
			continue
		}
		if start >= 0 && unicode.IsUpper(r) {
			prev := rs[i-1]
			nextLower := i+1 < len(rs) && unicode.IsLower(rs[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || unicode.IsUpper(prev) && nextLower {
				flush(i)
			}
		}
		if start < 0 {
			start = i
		}
	}
	flush(len(rs))
	return words
}

// SnakeCase returns the words of s in lower case joined by underscores:
// "HTTP server error" becomes "http_server_error".
func SnakeCase(s string) string {
	words := Words(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// PascalCase returns the words of s capitalised and joined with no
// separator: "http server error" becomes "HttpServerError".
func PascalCase(s string) string {
	words := Words(s)
	for i, w := range words {
		words[i] = capitalize(w)
	}
	return strings.Join(words, "")
}

// CapitalCase returns every word of s capitalised and joined by spaces:
// "the_lord_of_the_rings" becomes "The Lord Of The Rings".
func CapitalCase(s string) string {
	words := Words(s)
	for i, w := range words {
		words[i] = capitalize(w)
	}
	return strings.Join(words, " ")
}

// minorWords stay lower case in [TitleCase] unless first or last.
var minorWords = map[string]bool{
	"a": true, "an": true, "the": true,
	"and": true, "but": true, "or": true, "nor": true, "for": true, "so": true, "yet": true,
	"as": true, "at": true, "by": true, "in": true, "of": true, "off": true,
	"on": true, "per": true, "to": true, "up": true, "via": true,
}

// TitleCase is [CapitalCase] that keeps articles, short conjunctions and
// short prepositions in lower case except as the first or last word:
// "the lord of the rings" becomes "The Lord of the Rings".
func TitleCase(s string) string {
	words := Words(s)
	for i, w := range words {
		lower := strings.ToLower(w)
		if minorWords[lower] && i != 0 && i != len(words)-1 {
			words[i] = lower
		} else {
			words[i] = capitalize(w)
		}
	}
	return strings.Join(words, " ")
}

// capitalize upper-cases the first rune of w and lower-cases the rest.
func capitalize(w string) string {
	rs := []rune(strings.ToLower(w))
	if len(rs) > 0 {
		rs[0] = unicode.ToTitle(rs[0])
	}
	return string(rs)
}
