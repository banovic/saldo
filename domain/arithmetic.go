package domain

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrOverflow = errors.New("overflow")
)

// This is part of low level infrastructure - language level utility.

// mul returns the product of all numbers in ns.
// If there is an overflow it returns an error.
// With no arguments it returns 1 (multiplicative identity).
func mul(ns ...int64) (int64, error) {
	p := big.NewInt(1)
	for _, n := range ns {
		p.Mul(p, big.NewInt(n))
	}
	if !p.IsInt64() {
		return 0, fmt.Errorf("%w: mul: %v cannot be represented by int64", ErrOverflow, p)
	}
	return p.Int64(), nil
}

// add returns the sum of all the numbers in ns.
// If there is an overflow it returns an error.
// With no arguments it returns 0 (addition identity).
func add(ns ...int64) (int64, error) {
	s := big.NewInt(0)
	for _, n := range ns {
		s.Add(s, big.NewInt(n))
	}
	if !s.IsInt64() {
		return 0, fmt.Errorf("%w: add: %v cannot be represented by int64", ErrOverflow, s)
	}
	return s.Int64(), nil
}

// divAndRoundHalfUp accepts 2 integers (int64), and calculates division using the half-up rounding.
// Half up rounding means that if the result of division is middle it goes away from zero, this
// applies to both - positive and negative numbers.
// If den <= 0 it returns an error.
// For example:
//
//	 7 / 2 =  3.5 ->  4
//	-7 / 2 = -3.5 -> -4
//	 7 / 4 =  1.75 -> 2
//	 7 / 5 =  1.2  -> 1
func divAndRoundHalfUp(num, den int64) (int64, error) {
	if den <= 0 {
		return 0, fmt.Errorf("divAndRoundHalfUp: den must be greater than zero, got %d", den)
	}
	// Go truncates toward zero.
	p, r := num/den, num%den
	if r == 0 {
		return p, nil
	}
	if r < 0 {
		r = -r
	}
	if r < den-r {
		return p, nil
	}
	if num < 0 {
		return p - 1, nil
	}
	return p + 1, nil
}
