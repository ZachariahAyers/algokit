// Package numeric implements integer algorithms: factorials, Fibonacci
// numbers, primes, powers, roots, Roman numerals, English number names and
// a few classic exercises.
//
// Functions that can overflow uint64 or int64 return [ErrOverflow] rather
// than a wrapped result.
package numeric

import "errors"

// ErrOverflow reports a result too large for the return type.
var ErrOverflow = errors.New("numeric: result overflows")

// ErrDomain reports an argument outside the function's domain.
var ErrDomain = errors.New("numeric: argument out of domain")
