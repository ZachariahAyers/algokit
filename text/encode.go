package text

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrNoMorse reports a rune with no International Morse Code encoding.
var ErrNoMorse = errors.New("text: no Morse encoding")

var morse = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
}

// ToMorse encodes the letters and digits of s in International Morse Code,
// ignoring case. Letters are separated by a space and words by " / ".
//
// Time O(n). Space O(n). Returns [ErrNoMorse] for any other rune that is
// not whitespace.
func ToMorse(s string) (string, error) {
	words := strings.Fields(s)
	encoded := make([]string, len(words))
	for i, w := range words {
		letters := make([]string, 0, len(w))
		for _, r := range w {
			code, ok := morse[unicode.ToUpper(r)]
			if !ok {
				return "", fmt.Errorf("%w: %q", ErrNoMorse, r)
			}
			letters = append(letters, code)
		}
		encoded[i] = strings.Join(letters, " ")
	}
	return strings.Join(encoded, " / "), nil
}

// Caesar shifts each ASCII letter of s by shift places through the
// alphabet, wrapping around and keeping case. Other runes are unchanged.
// A negative shift moves backwards, so Caesar(Caesar(s, k), -k) == s.
//
// Time O(n). Space O(n).
func Caesar(s string, shift int) string {
	k := rune(((shift % 26) + 26) % 26)
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+k)%26
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+k)%26
		}
		return r
	}, s)
}

// Polynomial hash parameters: a base above every byte value and a prime
// modulus small enough that products stay inside uint64.
const (
	HashBase    = 257
	HashModulus = 1_000_000_007
)

// PolynomialHash returns the polynomial hash of the bytes of s:
// (s[0]·B^(n-1) + s[1]·B^(n-2) + … + s[n-1]) mod M, with B = [HashBase]
// and M = [HashModulus]. This form lets a rolling hash (as in Rabin-Karp)
// slide a window one byte at a time. Equal strings hash equally; different
// strings collide with probability about 1/M.
//
// Time O(n). Space O(1).
func PolynomialHash(s string) uint64 {
	var h uint64
	for i := range len(s) {
		h = (h*HashBase + uint64(s[i])) % HashModulus
	}
	return h
}
