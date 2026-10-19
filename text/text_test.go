package text

import (
	"errors"
	"maps"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	yes := []string{"", "a", "racecar", "A man, a plan, a canal: Panama", "No 'x' in Nixon", "été", "12321"}
	no := []string{"ab", "racecars", "0P", "été!x"}
	for _, s := range yes {
		if !IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = false", s)
		}
	}
	for _, s := range no {
		if IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = true", s)
		}
	}
}

func bruteLongestPalindrome(s string) string {
	rs := []rune(s)
	best := ""
	for i := range rs {
		for j := i; j <= len(rs); j++ {
			sub := rs[i:j]
			rev := slices.Clone(sub)
			slices.Reverse(rev)
			if slices.Equal(sub, rev) && len(sub) > len([]rune(best)) {
				best = string(sub)
			}
		}
	}
	return best
}

func TestLongestPalindrome(t *testing.T) {
	cases := map[string]string{
		"":                 "",
		"x":                "x",
		"ab":               "a",
		"babad":            "bab",
		"cbbd":             "bb",
		"forgeeksskeegfor": "geeksskeeg",
		"añaña":            "añaña",
	}
	for in, want := range cases {
		if got := LongestPalindrome(in); got != want {
			t.Errorf("LongestPalindrome(%q) = %q, want %q", in, got, want)
		}
	}
	rng := rand.New(rand.NewPCG(11, 12))
	for range 500 {
		rs := make([]rune, rng.IntN(25))
		for i := range rs {
			rs[i] = []rune("abé")[rng.IntN(3)]
		}
		s := string(rs)
		if got, want := LongestPalindrome(s), bruteLongestPalindrome(s); got != want {
			t.Fatalf("LongestPalindrome(%q) = %q, want %q", s, got, want)
		}
	}
}

func TestAnagrams(t *testing.T) {
	if got := Anagrams("cab"); !slices.Equal(got, []string{"abc", "acb", "bac", "bca", "cab", "cba"}) {
		t.Fatalf("got %q", got)
	}
	if got := Anagrams("aab"); !slices.Equal(got, []string{"aab", "aba", "baa"}) {
		t.Fatalf("repeats: got %q", got)
	}
	if got := Anagrams(""); !slices.Equal(got, []string{""}) {
		t.Fatalf("empty: got %q", got)
	}
}

func TestCharFrequency(t *testing.T) {
	got := CharFrequency("héllo")
	want := map[rune]int{'h': 1, 'é': 1, 'l': 2, 'o': 1}
	if !maps.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
	if len(CharFrequency("")) != 0 {
		t.Fatal("empty")
	}
}

func TestReverseAndRotation(t *testing.T) {
	if Reverse("héllo") != "olléh" || Reverse("") != "" {
		t.Fatal("Reverse")
	}
	if got := ReverseWords("  the sky   is blue "); got != "blue is sky the" {
		t.Fatalf("ReverseWords = %q", got)
	}
	if !IsRotation("waterbottle", "erbottlewat") || !IsRotation("", "") {
		t.Fatal("IsRotation true cases")
	}
	if IsRotation("abc", "acb") || IsRotation("ab", "aba") {
		t.Fatal("IsRotation false cases")
	}
}

func TestRemoveVowels(t *testing.T) {
	if got := RemoveVowels("Education is Optional"); got != "dctn s ptnl" {
		t.Fatalf("got %q", got)
	}
}

func TestWords(t *testing.T) {
	cases := map[string][]string{
		"":                 nil,
		"hello world":      {"hello", "world"},
		"fooBar":           {"foo", "Bar"},
		"HTTPServerError":  {"HTTP", "Server", "Error"},
		"snake_case-kebab": {"snake", "case", "kebab"},
		"v2Api":            {"v2", "Api"},
		"  --  ":           nil,
		"ÉcoleNormale":     {"École", "Normale"},
	}
	for in, want := range cases {
		if got := Words(in); !slices.Equal(got, want) {
			t.Errorf("Words(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCaseConversions(t *testing.T) {
	type conv struct {
		name string
		fn   func(string) string
		in   string
		want string
	}
	cases := []conv{
		{"Snake", SnakeCase, "HTTP server Error", "http_server_error"},
		{"Snake", SnakeCase, "fooBarBaz", "foo_bar_baz"},
		{"Pascal", PascalCase, "http server error", "HttpServerError"},
		{"Pascal", PascalCase, "snake_case", "SnakeCase"},
		{"Capital", CapitalCase, "the_lord_of_the_rings", "The Lord Of The Rings"},
		{"Title", TitleCase, "the lord of the rings", "The Lord of the Rings"},
		{"Title", TitleCase, "what are you looking at", "What Are You Looking At"},
		{"Title", TitleCase, "", ""},
	}
	for _, c := range cases {
		if got := c.fn(c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
	}
}

func TestToMorse(t *testing.T) {
	got, err := ToMorse("SOS sos 73")
	if err != nil || got != "... --- ... / ... --- ... / --... ...--" {
		t.Fatalf("got %q, %v", got, err)
	}
	if got, err := ToMorse("   "); err != nil || got != "" {
		t.Fatalf("blank: %q, %v", got, err)
	}
	if _, err := ToMorse("hi!"); !errors.Is(err, ErrNoMorse) {
		t.Fatalf("error = %v, want ErrNoMorse", err)
	}
}

func TestCaesar(t *testing.T) {
	if got := Caesar("Hello, World!", 3); got != "Khoor, Zruog!" {
		t.Fatalf("got %q", got)
	}
	if got := Caesar("xyz XYZ", 3); got != "abc ABC" {
		t.Fatalf("wrap: got %q", got)
	}
	for _, k := range []int{0, 1, 25, 26, 27, -1, -53, 1000} {
		s := "The quick brown fox, é."
		if back := Caesar(Caesar(s, k), -k); back != s {
			t.Fatalf("shift %d round trip: %q", k, back)
		}
	}
	if Caesar("abc", 26) != "abc" {
		t.Fatal("shift 26 should be identity")
	}
}

func TestPolynomialHash(t *testing.T) {
	if PolynomialHash("") != 0 {
		t.Fatal("empty hash")
	}
	if got := PolynomialHash("ab"); got != 'a'*HashBase+'b' {
		t.Fatalf("got %d", got)
	}
	// Rolling property: drop the first byte and append one in O(1).
	s := "abracadabra"
	const w = 4
	pow := uint64(1)
	for range w - 1 {
		pow = pow * HashBase % HashModulus
	}
	h := PolynomialHash(s[:w])
	for i := 1; i+w <= len(s); i++ {
		h = (h + HashModulus - uint64(s[i-1])*pow%HashModulus) % HashModulus
		h = (h*HashBase + uint64(s[i+w-1])) % HashModulus
		if want := PolynomialHash(s[i : i+w]); h != want {
			t.Fatalf("window %d: rolled %d, direct %d", i, h, want)
		}
	}
	seen := map[uint64]string{}
	for _, a := range "abcdefghij" {
		for _, b := range "abcdefghij" {
			for _, c := range "abcdefghij" {
				str := string([]rune{a, b, c})
				if prev, ok := seen[PolynomialHash(str)]; ok {
					t.Fatalf("collision %q %q", prev, str)
				}
				seen[PolynomialHash(str)] = str
			}
		}
	}
}
