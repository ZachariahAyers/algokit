package numeric

import (
	"errors"
	"math"
	"slices"
	"testing"
)

func TestFactorial(t *testing.T) {
	want := []uint64{1, 1, 2, 6, 24, 120, 720}
	for n, w := range want {
		if got, err := Factorial(uint(n)); err != nil || got != w {
			t.Fatalf("Factorial(%d) = %d, %v; want %d", n, got, err, w)
		}
	}
	if got, err := Factorial(20); err != nil || got != 2432902008176640000 {
		t.Fatalf("Factorial(20) = %d, %v", got, err)
	}
	if _, err := Factorial(21); !errors.Is(err, ErrOverflow) {
		t.Fatalf("Factorial(21) error = %v, want ErrOverflow", err)
	}
}

func TestFibonacci(t *testing.T) {
	want := []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	for n, w := range want {
		if got, err := Fibonacci(uint(n)); err != nil || got != w {
			t.Fatalf("Fibonacci(%d) = %d, %v; want %d", n, got, err, w)
		}
	}
	if got, err := Fibonacci(93); err != nil || got != 12200160415121876738 {
		t.Fatalf("Fibonacci(93) = %d, %v", got, err)
	}
	if _, err := Fibonacci(94); !errors.Is(err, ErrOverflow) {
		t.Fatalf("Fibonacci(94) error = %v, want ErrOverflow", err)
	}
}

func TestCollatzSteps(t *testing.T) {
	cases := map[uint64]int{1: 0, 2: 1, 6: 8, 7: 16, 27: 111}
	for n, want := range cases {
		if got, err := CollatzSteps(n); err != nil || got != want {
			t.Fatalf("CollatzSteps(%d) = %d, %v; want %d", n, got, err, want)
		}
	}
	if _, err := CollatzSteps(0); !errors.Is(err, ErrDomain) {
		t.Fatalf("CollatzSteps(0) error = %v, want ErrDomain", err)
	}
	if _, err := CollatzSteps(math.MaxUint64); !errors.Is(err, ErrOverflow) {
		t.Fatalf("CollatzSteps(max) error = %v, want ErrOverflow", err)
	}
}

func TestFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %q", got)
	}
	if got := FizzBuzz(-3); len(got) != 0 {
		t.Fatalf("FizzBuzz(-3) = %q, want empty", got)
	}
}

func TestIsPrimeAgreesWithSieve(t *testing.T) {
	primes := Sieve(10_000)
	if !slices.Equal(primes[:10], []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}) {
		t.Fatalf("first primes: %v", primes[:10])
	}
	if len(primes) != 1229 {
		t.Fatalf("pi(10000) = %d, want 1229", len(primes))
	}
	for n := range 10_001 {
		_, isPrime := slices.BinarySearch(primes, n)
		if IsPrime(uint64(n)) != isPrime {
			t.Fatalf("IsPrime(%d) = %v", n, !isPrime)
		}
	}
	// The largest prime below 2^32, and a semiprime of two primes near 2^16.
	if !IsPrime(4294967291) {
		t.Fatal("2^32-5 is prime")
	}
	if IsPrime(65521 * 65519) {
		t.Fatal("product of two primes reported prime")
	}
	if Sieve(1) != nil || !slices.Equal(Sieve(2), []int{2}) {
		t.Fatal("Sieve small bounds")
	}
}

func TestPow(t *testing.T) {
	cases := []struct{ base, exp, want uint64 }{
		{0, 0, 1}, {0, 5, 0}, {1, 1 << 40, 1}, {2, 10, 1024}, {3, 40, 12157665459056928801}, {2, 63, 1 << 63},
	}
	for _, c := range cases {
		if got, err := Pow(c.base, c.exp); err != nil || got != c.want {
			t.Fatalf("Pow(%d, %d) = %d, %v; want %d", c.base, c.exp, got, err, c.want)
		}
	}
	for _, c := range [][2]uint64{{2, 64}, {3, 41}, {1 << 32, 2}} {
		if _, err := Pow(c[0], c[1]); !errors.Is(err, ErrOverflow) {
			t.Fatalf("Pow(%d, %d) error = %v, want ErrOverflow", c[0], c[1], err)
		}
	}
}

func TestIsPowerOf(t *testing.T) {
	yes := [][2]uint64{{1, 2}, {8, 2}, {81, 3}, {1 << 63, 2}, {0, 0}, {1, 0}, {1, 1}, {1000, 10}}
	no := [][2]uint64{{0, 2}, {6, 2}, {2, 0}, {2, 1}, {12, 3}, {9, 2}}
	for _, c := range yes {
		if !IsPowerOf(c[0], c[1]) {
			t.Fatalf("IsPowerOf(%d, %d) = false", c[0], c[1])
		}
	}
	for _, c := range no {
		if IsPowerOf(c[0], c[1]) {
			t.Fatalf("IsPowerOf(%d, %d) = true", c[0], c[1])
		}
	}
}

func TestSqrt(t *testing.T) {
	for n := range uint64(20_000) {
		r := Sqrt(n)
		if r*r > n || (r+1)*(r+1) <= n {
			t.Fatalf("Sqrt(%d) = %d", n, r)
		}
	}
	for _, n := range []uint64{math.MaxUint64, 1<<62 - 1, 1 << 62, 4294967295 * 4294967295} {
		r := Sqrt(n)
		hiOK := r <= math.MaxUint32
		if !hiOK || r*r > n || (r+1 <= math.MaxUint32 && (r+1)*(r+1) <= n) {
			t.Fatalf("Sqrt(%d) = %d", n, r)
		}
	}
}

func TestReverseDigits(t *testing.T) {
	cases := map[int64]int64{0: 0, 7: 7, 120: 21, -120: -21, 1534236469: 9646324351}
	for n, want := range cases {
		if got, err := ReverseDigits(n); err != nil || got != want {
			t.Fatalf("ReverseDigits(%d) = %d, %v; want %d", n, got, err, want)
		}
	}
	// Reversals that land exactly on the int64 limits still fit.
	if got, err := ReverseDigits(7085774586302733229); err != nil || got != 9223372036854775807 {
		t.Fatalf("ReverseDigits to MaxInt64 = %d, %v", got, err)
	}
	if got, err := ReverseDigits(-8085774586302733229); err != nil || got != math.MinInt64 {
		t.Fatalf("ReverseDigits to MinInt64 = %d, %v", got, err)
	}
	for _, n := range []int64{1999999999999999999, -1999999999999999999, 8085774586302733229} {
		if _, err := ReverseDigits(n); !errors.Is(err, ErrOverflow) {
			t.Fatalf("ReverseDigits(%d) error = %v, want ErrOverflow", n, err)
		}
	}
}

func TestSameDigits(t *testing.T) {
	if !SameDigits(1123, 3121) || !SameDigits(-45, 54) || !SameDigits(0, 0) {
		t.Fatal("expected matches")
	}
	if SameDigits(12, 122) || SameDigits(10, 1) || SameDigits(0, 10) {
		t.Fatal("expected mismatches")
	}
}

func TestRomanRoundTrip(t *testing.T) {
	for n := 1; n <= 3999; n++ {
		s, err := ToRoman(n)
		if err != nil {
			t.Fatalf("ToRoman(%d): %v", n, err)
		}
		back, err := FromRoman(s)
		if err != nil || back != n {
			t.Fatalf("FromRoman(%q) = %d, %v; want %d", s, back, err, n)
		}
	}
	known := map[int]string{4: "IV", 9: "IX", 14: "XIV", 40: "XL", 1994: "MCMXCIV", 2026: "MMXXVI", 3999: "MMMCMXCIX"}
	for n, want := range known {
		if got, _ := ToRoman(n); got != want {
			t.Fatalf("ToRoman(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestRomanRejectsInvalid(t *testing.T) {
	for _, n := range []int{0, -1, 4000} {
		if _, err := ToRoman(n); !errors.Is(err, ErrDomain) {
			t.Fatalf("ToRoman(%d) error = %v", n, err)
		}
	}
	for _, s := range []string{"", "IIII", "VV", "IC", "MMMM", "iv", "XLX", "ABC", "IIX"} {
		if _, err := FromRoman(s); !errors.Is(err, ErrDomain) {
			t.Fatalf("FromRoman(%q) error = %v", s, err)
		}
	}
}

func TestToWords(t *testing.T) {
	cases := map[int64]string{
		0:          "zero",
		7:          "seven",
		13:         "thirteen",
		20:         "twenty",
		42:         "forty-two",
		100:        "one hundred",
		101:        "one hundred one",
		1_000:      "one thousand",
		1_234:      "one thousand two hundred thirty-four",
		1_000_001:  "one million one",
		-5:         "minus five",
		90_000_019: "ninety million nineteen",
		math.MinInt64: "minus nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion " +
			"thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight",
	}
	for n, want := range cases {
		if got := ToWords(n); got != want {
			t.Fatalf("ToWords(%d) = %q, want %q", n, got, want)
		}
	}
}
