// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import "math"

const (
	// macheps is the machine epsilon aka unit roundoff
	// The machine epsilon is an upper bound on the absolute relative true error in
	// representing a number.
	// If y is the machine representation of x then |(x-y)/x| <= macheps
	// https://en.wikipedia.org/wiki/Machine_epsilon
	// Golangs float64 type has a 52-bit fractional mantissa,
	// therefore the value 2^-52
	macheps = 1.0 / (1 << 52)
)

// float64Pred dynamically updates the potential error
type float64Pred struct {
	// n is the number
	n float64
	// e is the max rounding error possible
	e float64
}

// newFloat64Pred returns a new float64Pred e set to 0.
func newFloat64Pred(n float64) float64Pred {
	return float64Pred{
		n: n,
		e: 0,
	}
}

// addFloat64 adds f to p and updates the potential error
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
// e = macheps*|x|.
func (p float64Pred) addFloat64(f float64) float64Pred {
	p.n += f
	p.e += macheps * math.Abs(p.n)
	return p
}

// addFloat64Pred adds b to a and updates the potential error
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
// e = macheps*|x|.
func (a float64Pred) addFloat64Pred(b float64Pred) float64Pred {
	a.n += b.n
	a.e += macheps*math.Abs(a.n) + b.e
	return a
}

// subFloat64 subtracts f from p and updates the potential error
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
// e = macheps*|x|.
func (p float64Pred) subFloat64(f float64) float64Pred {
	p.n -= f
	p.e += macheps * math.Abs(p.n)
	return p
}

// subFloat64Pred subtracts f from p and updates the potential error
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
// e = macheps*|x|.
func (a float64Pred) subFloat64Pred(b float64Pred) float64Pred {
	a.n -= b.n
	a.e += macheps*math.Abs(a.n) + b.e
	return a
}

// mulFloat64 multiplies p with f and updates the potential error.
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
// e = macheps*|x|.
//
// mul(mul(a,b),c) = mul(a*b+error,c) = a*b*c + error*c + newError
// sum(mul(a,b),c) = sum(a*b+error,c) = a*b+c + error + newError
// Conclusively, when multiplications are chained, the error also depends on the value
// of the number, but this does not apply to sums or subtractions.
//
//If this is not a chained multiplication p.e will be 0, making that part irrelevant.
func (p float64Pred) mulFloat64(f float64) float64Pred {
	p.n *= f
	p.e += macheps*math.Abs(p.n) + p.e*f
	return p
}
