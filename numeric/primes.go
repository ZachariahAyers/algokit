package numeric

// IsPrime reports whether n is prime, by trial division over 6k ± 1.
//
// Time O(√n). Space O(1). Near 2^64 that is about 7×10^8 divisions; for
// large inputs prefer math/big's ProbablyPrime.
func IsPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	if n < 4 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	// i <= n/i avoids the overflow i*i would hit near 2^64.
	for i := uint64(5); i <= n/i; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// Sieve returns every prime <= n in ascending order, using the sieve of
// Eratosthenes.
//
// Time O(n log log n). Space O(n).
func Sieve(n int) []int {
	if n < 2 {
		return nil
	}
	composite := make([]bool, n+1)
	var primes []int
	for i := 2; i <= n; i++ {
		if composite[i] {
			continue
		}
		primes = append(primes, i)
		for j := i * i; j <= n; j += i {
			composite[j] = true
		}
	}
	return primes
}
