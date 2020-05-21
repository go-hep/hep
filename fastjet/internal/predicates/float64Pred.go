// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"math"
)

const (
	// macheps is the machine epsilon aka unit roundoff
	// The machine epsilon is an upper bound on the absolute relative true error in
	// representing a number.
	// If y is the machine representation of x then |(x-y)/x| <= macheps
	// https://en.wikipedia.org/wiki/Machine_epsilon
	// Go's float64 type has a 52-bit fractional mantissa,
	// therefore the value 2^-52
	macheps = 1.0 / (1 << 52)
)

// float64Pred dynamically updates the potential error.
//
// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e.
// Since we want the max possible error we assume |(x-y)/x| = macheps
// macheps*|x| = |x - y|
// if (x-y)>=0 then macheps*|x| = x - y  ->  y = x- macheps*|x|
// else macheps*|x| = -(x - y)  ->  macheps*|x|+x = y
// Each one of these has two cases again. x can be positive or negative, resulting in the same two possible equations:
// y/(1-macheps)=x or y/(1+macheps)=x.
// Since x is unknown, we will use the larger number as factor, to avoid having an error greater than
// the maxError value we have.
// |y-x| = macheps*|x| -> e = macheps*|x| -> e = macheps* |y|/(1-macheps)
// A Special case is when y = 0. Then we use the smallest nonzero float, because that is the max
// possible error in this case.
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

// // addFloat64 adds f to p and updates the potential error
// func (p float64Pred) addFloat64(f float64) float64Pred {
// 	p.n += f
// 	if p.n == 0 {
// 		p.e += math.SmallestNonzeroFloat64
// 	} else {
// 		p.e += macheps * math.Abs(p.n) / (1 - macheps)
// 	}
// 	return p
// }

// addFloat64Pred adds b to a and updates the potential error
func (a float64Pred) addFloat64Pred(b float64Pred) float64Pred {
	a.n += b.n
	if a.n == 0 {
		a.e += math.SmallestNonzeroFloat64 + b.e
	} else {
		a.e += macheps*math.Abs(a.n)/(1-macheps) + b.e
	}
	return a
}

// // subFloat64 subtracts f from p and updates the potential error
// func (p float64Pred) subFloat64(f float64) float64Pred {
// 	p.n -= f
// 	if p.n == 0 {
// 		p.e += math.SmallestNonzeroFloat64
// 	} else {
// 		p.e += macheps * math.Abs(p.n) / (1 - macheps)
// 	}
// 	return p
// }

// subFloat64Pred subtracts f from p and updates the potential error
func (a float64Pred) subFloat64Pred(b float64Pred) float64Pred {
	a.n -= b.n
	if a.n == 0 {
		a.e += math.SmallestNonzeroFloat64 + b.e
	} else {
		a.e += macheps*math.Abs(a.n)/(1-macheps) + b.e
	}
	return a
}

// mulFloat64 multiplies p with f and updates the potential error.
//
// mul(mul(a,b),c) = mul(a*b+error,c) = a*b*c + error*c + newError
// sum(mul(a,b),c) = sum(a*b+error,c) = a*b+c + error + newError
// Conclusively, when multiplications are chained, the error also depends on the value
// of the number, but this does not apply to sums or subtractions.
//
//If this is not a chained multiplication p.e will be 0, making that part irrelevant.
func (p float64Pred) mulFloat64(f float64) float64Pred {
	p.n *= f
	if p.n == 0 {
		p.e += math.SmallestNonzeroFloat64 + p.e*math.Abs(f)
	} else {
		p.e += macheps*math.Abs(p.n)/(1-macheps) + p.e*math.Abs(f)
	}
	return p
}
