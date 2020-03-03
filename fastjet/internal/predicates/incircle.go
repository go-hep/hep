// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/mat"
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
}

// Incircle determines the relative position of the point (x,y) in relation to the circle formed
// by the three points (x1,y1),(x2,y2) and (x3,y3). The three points have to be ordered counterclockwise or
// Outside and Inside will be reversed.
func Incircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	pos := simpleIncircle(x1, y1, x2, y2, x3, y3, x, y)
	if pos == IndeterminatePosition {
		pos = matIncircle(x1, y1, x2, y2, x3, y3, x, y)
	}
	return pos
}

// simpleIncircle determines the relative position using the simple float64 type.
// Its accuracy can't be guaranteed, therefore close decisions
// return IndeterminatePosition which signals Incircle that further
// testing is necessary.
//
// It computes the determinant of the matrix and returns the relative position
// based on the value of the determinant.
// |x1 y1 x1^2+y1^2 1|
// |x2 y2 x2^2+y2^2 1|
// |x3 y3 x3^2+y3^2 1|
// |x  y  x^2 +y^2  1|
func simpleIncircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	m := []float64{
		x1, y1, x1*x1 + y1*y1, 1,
		x2, y2, x2*x2 + y2*y2, 1,
		x3, y3, x3*x3 + y3*y3, 1,
		x, y, x*x + y*y, 1,
	}
	p := rowFloat(3, 2, 1, 0, false, m).addFloat64Pred(rowFloat(3, 1, 2, 0, true, m)).addFloat64Pred(
		rowFloat(2, 1, 3, 0, false, m)).addFloat64Pred(rowFloat(3, 2, 0, 1, true, m)).addFloat64Pred(
		rowFloat(3, 0, 2, 1, false, m)).addFloat64Pred(rowFloat(2, 0, 3, 1, true, m)).addFloat64Pred(
		rowFloat(3, 1, 0, 2, false, m)).addFloat64Pred(rowFloat(3, 0, 1, 2, true, m)).addFloat64Pred(
		rowFloat(1, 0, 3, 2, false, m)).addFloat64Pred(rowFloat(2, 1, 0, 3, true, m)).addFloat64Pred(
		rowFloat(2, 0, 1, 3, false, m)).addFloat64Pred(rowFloat(1, 0, 2, 3, true, m))
	// det := m[0][3]*m[1][2]*m[2][1]*m[3][0] - m[0][2]*m[1][3]*m[2][1]*m[3][0] -
	//	m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][1]*m[1][3]*m[2][2]*m[3][0] +
	//	m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][1]*m[1][2]*m[2][3]*m[3][0] -
	//	m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][2]*m[1][3]*m[2][0]*m[3][1] +
	//	m[0][3]*m[1][0]*m[2][2]*m[3][1] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
	//	m[0][2]*m[1][0]*m[2][3]*m[3][1] + m[0][0]*m[1][2]*m[2][3]*m[3][1] +
	//	m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][1]*m[1][3]*m[2][0]*m[3][2] -
	//	m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][0]*m[1][3]*m[2][1]*m[3][2] +
	//	m[0][1]*m[1][0]*m[2][3]*m[3][2] - m[0][0]*m[1][1]*m[2][3]*m[3][2] -
	//	m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][1]*m[1][2]*m[2][0]*m[3][3] +
	//	m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][0]*m[1][2]*m[2][1]*m[3][3] -
	//	m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][0]*m[1][1]*m[2][2]*m[3][3]
	det := p.n
	// e determines when the determinant in simpleIncircle is too close to 0 to rely on floating point operations.
	e := p.e
	if det < -e {
		return Outside
	}
	if det > e {
		return Inside
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

// robustIncircle determines the relative position using the accurate big/Rat type.
//
// It computes the determinant of the matrix and returns the relative position
// based on the value of the determinant.
// |x1 y1 x1^2+y1^2 1|
// |x2 y2 x2^2+y2^2 1|
// |x3 y3 x3^2+y3^2 1|
// |x  y  x^2 +y^2  1|
func robustIncircle(x1, y1, x2, y2, x3, y3, x, y *big.Rat) RelativePosition {
	m := []*big.Rat{
		x1, y1, bigAdd(bigMul(x1, x1), bigMul(y1, y1)), one,
		x2, y2, bigAdd(bigMul(x2, x2), bigMul(y2, y2)), one,
		x3, y3, bigAdd(bigMul(x3, x3), bigMul(y3, y3)), one,
		x, y, bigAdd(bigMul(x, x), bigMul(y, y)), one,
	}
	// det := m[0][3]*m[1][2]*m[2][1]*m[3][0] - m[0][2]*m[1][3]*m[2][1]*m[3][0] -
	//	m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][1]*m[1][3]*m[2][2]*m[3][0] +
	//	m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][1]*m[1][2]*m[2][3]*m[3][0] -
	//	m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][2]*m[1][3]*m[2][0]*m[3][1] +
	//	m[0][3]*m[1][0]*m[2][2]*m[3][1] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
	//	m[0][2]*m[1][0]*m[2][3]*m[3][1] + m[0][0]*m[1][2]*m[2][3]*m[3][1] +
	//	m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][1]*m[1][3]*m[2][0]*m[3][2] -
	//	m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][0]*m[1][3]*m[2][1]*m[3][2] +
	//	m[0][1]*m[1][0]*m[2][3]*m[3][2] - m[0][0]*m[1][1]*m[2][3]*m[3][2] -
	//	m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][1]*m[1][2]*m[2][0]*m[3][3] +
	//	m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][0]*m[1][2]*m[2][1]*m[3][3] -
	//	m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][0]*m[1][1]*m[2][2]*m[3][3]
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
	sign := det.Sign()
	switch sign {
	case 1:
		return Inside
	case -1:
		return Outside
	case 0:
		return On
	}
	return IndeterminatePosition
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

// matIncircle determines the relative position using the mat package.
//
// It first computes the conditional number of the matrix. When the condition number
// is higher than the Condition Tolerance, then we assume the matrix is singular and
// the determinant is 0. If the determinant is not 0 the sign of the determinant is computed.
// |x1 y1 x1^2+y1^2 1|
// |x2 y2 x2^2+y2^2 1|
// |x3 y3 x3^2+y3^2 1|
// |x  y  x^2 +y^2  1|
func matIncircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	m := mat.NewDense(4, 4, []float64{x1, y1, x1*x1 + y1*y1, 1, x2, y2, x2*x2 + y2*y2, 1, x3, y3, x3*x3 + y3*y3, 1, x, y, x*x + y*y, 1})
	var lu mat.LU
	lu.Factorize(m)
	cond := lu.Cond()
	if cond > mat.ConditionTolerance {
		return On
	}
	// Since only the sign is needed LogDet achieves the result in faster time.
	_, sign := lu.LogDet()
	switch sign {
	case 1:
		return Inside
	case -1:
		return Outside
	}
	return IndeterminatePosition
}
