// Package num holds the numeric type constraints shared by algokit packages.
package num

// Signed matches every signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned matches every unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer matches every integer type.
type Integer interface {
	Signed | Unsigned
}

// Float matches every floating-point type.
type Float interface {
	~float32 | ~float64
}

// Real matches every integer and floating-point type.
type Real interface {
	Integer | Float
}
