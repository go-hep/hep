// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"fmt"
	"math/big"
)

// RelativePosition is the position of a point relative to a circle
type RelativePosition int

const (
	// Inside the circle
	Inside RelativePosition = iota
	// On the circle
	On
	// Outside the circle
	Outside
	IndeterminatePosition
)

func (p RelativePosition) String() string {
	switch p {
	case Inside:
		return "Inside Circle"
	case On:
		return "On Circle"
	case Outside:
		return "Outside Circle"
	case IndeterminatePosition:
		return "Indeterminate"
	default:
		panic(fmt.Errorf("predicates: unknown RelativePosition %d", int(p)))
	}
	panic("unreachable")
}

// Incircle determines the relative position of the point (x,y) in relation to the circle formed
// by the three points (x1,y1),(x2,y2) and (x3,y3)
func Incircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	pos := simpleIncircle(x1, y1, x2, y2, x3, y3, x, y)
	if pos == IndeterminatePosition {
		pos = robustIncircle(setBig(x1), setBig(y1), setBig(x2), setBig(y2), setBig(x3), setBig(y3), setBig(x), setBig(y))
	}
	return pos
}

// simpleIncircle finds the determinant of the matrix using the simple float64 type.
// Then it returns the relative position based on the value of the determinant.
// |1 x1 y1 x1^2+y1^2|
// |1 x2 y2 x2^2+y2^2|
// |1 x3 y3 x3^2+y3^2|
// |1 x  y  x^2 +y^2 |
func simpleIncircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	m := []float64{
		1, x1, y1, x1*x1 + y1*y1,
		1, x2, y2, x2*x2 + y2*y2,
		1, x3, y3, x3*x3 + y3*y3,
		1, x, y, x*x + y*y,
	}
	p := rowFloat(3, 2, 1, 0, false, m).addFloat64Pred(rowFloat(3, 1, 2, 0, true, m)).addFloat64Pred(
		rowFloat(2, 1, 3, 0, false, m)).addFloat64Pred(rowFloat(3, 2, 0, 1, true, m)).addFloat64Pred(
		rowFloat(3, 0, 2, 1, false, m)).addFloat64Pred(rowFloat(2, 0, 3, 1, true, m)).addFloat64Pred(
		rowFloat(3, 1, 0, 2, false, m)).addFloat64Pred(rowFloat(3, 0, 1, 2, true, m)).addFloat64Pred(
		rowFloat(1, 0, 3, 2, false, m)).addFloat64Pred(rowFloat(2, 1, 0, 3, true, m)).addFloat64Pred(
		rowFloat(2, 0, 1, 3, false, m)).addFloat64Pred(rowFloat(1, 0, 2, 3, true, m))
	/* det := m[3]*m[6]*m[9]*m[12] - m[2]*m[7]*m[9]*m[12] -
	m[3]*m[5]*m[10]*m[12] + m[1]*m[7]*m[10]*m[12] +
	m[2]*m[5]*m[11]*m[12] - m[1]*m[6]*m[11]*m[12] -
	m[3]*m[6]*m[8]*m[13] + m[2]*m[7]*m[8]*m[13] +
	m[3]*m[4]*m[10]*m[13] - m[0]*m[7]*m[10]*m[13] -
	m[2]*m[4]*m[11]*m[13] + m[0]*m[6]*m[11]*m[13] +
	m[3]*m[5]*m[8]*m[14] - m[1]*m[7]*m[8]*m[14] -
	m[3]*m[4]*m[9]*m[14] + m[0]*m[7]*m[9]*m[14] +
	m[1]*m[4]*m[11]*m[14] - m[0]*m[5]*m[11]*m[14] -
	m[2]*m[5]*m[8]*m[15] + m[1]*m[6]*m[8]*m[15] +
	m[2]*m[4]*m[9]*m[15] - m[0]*m[6]*m[9]*m[15] -
	m[1]*m[4]*m[10]*m[15] + m[0]*m[5]*m[10]*m[15] */
	det := p.n
	// e determines when the determinant in orientation is too close to 0 to rely on floating point operations.
	// Each intermediate result can have a potential absolute relative rounding error of macheps.
	// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
	// e = macheps*|x|.
	// mul(mul(a,b),c) = mul(a*b+macheps,c) = a*b*c + macheps*c + macheps
	// sum(mul(a,b),c) = sum(a*b+macheps,c) = a*b+c + macheps + macheps
	// Conclusively, when multiplications are chained, the error also depends on the value
	// of the number, but this does not apply to sums or subtractions.
	e := p.e
	if det < -e {
		return Inside
	}
	if det > e {
		return Outside
	}
	return IndeterminatePosition
}

// rowFloat is a helper function for robustIncircle
// If m[row][col] then each row in the determinant calculation is either
// m[0][a]*m[1][b]*m[2][c]*m[3][d] - m[0][b]*m[1][a]*m[2][c]*m[3][d] or
// - m[0][a]*m[1][b]*m[2][c]*m[3][d] + m[0][b]*m[1][a]*m[2][c]*m[3][d]
func rowFloat(a, b, c, d int, plus bool, m []float64) float64Pred {
	if plus {
		return newFloat64Pred(m[b]).mulFloat64(m[a+4]).mulFloat64(m[c+8]).mulFloat64(m[d+12]).subFloat64Pred(
			newFloat64Pred(m[a]).mulFloat64(m[b+4]).mulFloat64(m[c+8]).mulFloat64(m[d+12]))
	}
	return newFloat64Pred(m[a]).mulFloat64(m[b+4]).mulFloat64(m[c+8]).mulFloat64(m[d+12]).subFloat64Pred(
		newFloat64Pred(m[b]).mulFloat64(m[a+4]).mulFloat64(m[c+8]).mulFloat64(m[d+12]))
}

// robustIncircle computes the determinant of the matrix using the accurate big/Rat type.
// Then it returns the relative position based on the value of the determinant
// |1 x1 y1 x1^2+y1^2|
// |1 x2 y2 x2^2+y2^2|
// |1 x3 y3 x3^2+y3^2|
// |1 x  y  x^2 +y^2 |
func robustIncircle(x1, y1, x2, y2, x3, y3, x, y *big.Rat) RelativePosition {
	m := []*big.Rat{
		one, x1, y1, bigAdd(bigMul(x1, x1), bigMul(y1, y1)),
		one, x2, y2, bigAdd(bigMul(x2, x2), bigMul(y2, y2)),
		one, x3, y3, bigAdd(bigMul(x3, x3), bigMul(y3, y3)),
		one, x, y, bigAdd(bigMul(x, x), bigMul(y, y)),
	}
	det := bigAdd(
		bigAdd(
			bigAdd(
				bigAdd(rowBig(3, 2, 1, 0, false, m), rowBig(3, 1, 2, 0, true, m)),
				rowBig(2, 1, 3, 0, false, m),
			),
			bigAdd(
				bigAdd(rowBig(3, 2, 0, 1, true, m), rowBig(3, 0, 2, 1, false, m)),
				rowBig(2, 0, 3, 1, true, m),
			),
		),
		bigAdd(
			bigAdd(
				bigAdd(rowBig(3, 1, 0, 2, false, m), rowBig(3, 0, 1, 2, true, m)),
				rowBig(1, 0, 3, 2, false, m),
			),
			bigAdd(
				bigAdd(rowBig(2, 1, 0, 3, true, m), rowBig(2, 0, 1, 3, false, m)),
				rowBig(1, 0, 2, 3, true, m),
			),
		),
	)
	if det.Cmp(zero) < 0 {
		return Inside
	}
	if det.Cmp(zero) == 0 {
		return On
	}
	return Outside
}

// rowBig is a helper function for robustIncircle
// If m[row][col] then each row in the determinant calculation is either
// m[0][a]*m[1][b]*m[2][c]*m[3][d] - m[0][b]*m[1][a]*m[2][c]*m[3][d] or
// - m[0][a]*m[1][b]*m[2][c]*m[3][d] + m[0][b]*m[1][a]*m[2][c]*m[3][d]
func rowBig(a, b, c, d int, plus bool, m []*big.Rat) *big.Rat {
	if plus {
		return bigSub(
			bigMul(bigMul(m[b], m[a+4]), bigMul(m[c+8], m[d+12])),
			bigMul(bigMul(m[a], m[b+4]), bigMul(m[c+8], m[d+12])),
		)
	}
	return bigSub(
		bigMul(bigMul(m[a], m[b+4]), bigMul(m[c+8], m[d+12])),
		bigMul(bigMul(m[b], m[a+4]), bigMul(m[c+8], m[d+12])),
	)
}
