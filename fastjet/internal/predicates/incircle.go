// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"math/big"
)

const (
	// FIXME
	epsilonI = 1e-13
)

// Incircle determines whether the point (x,y) lies inside the circle formed
// by the three points (x1,y1),(x2,y2) and (x3,y3)
func Incircle(x1, y1, x2, y2, x3, y3, x, y float64) bool {
	det := simpleIncircle(x1, y1, x2, y2, x3, y3, x, y)
	if det > epsilonI {
		return false
	}
	if det < -epsilonI {
		return true
	}
	return robustIncircle(setBig(x1), setBig(y1), setBig(x2), setBig(y2), setBig(x3), setBig(y3), setBig(x), setBig(y))
}

// simpleIncircle finds the determinant of the matrix using the simple float64 type
// |1 x1 y1 x1^2+y1^2|
// |1 x2 y2 x2^2+y2^2|
// |1 x3 y3 x3^2+y3^2|
// |1 x  y  x^2 +y^2 |
func simpleIncircle(x1, y1, x2, y2, x3, y3, x, y float64) (det float64) {
	m := [][]float64{
		{1, x1, y1, x1*x1 + y1*y1},
		{1, x2, y2, x2*x2 + y2*y2},
		{1, x3, y3, x3*x3 + y3*y3},
		{1, x, y, x*x + y*y},
	}
	return m[0][3]*m[1][2]*m[2][1]*m[3][0] - m[0][2]*m[1][3]*m[2][1]*m[3][0] -
		m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][1]*m[1][3]*m[2][2]*m[3][0] +
		m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][1]*m[1][2]*m[2][3]*m[3][0] -
		m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][2]*m[1][3]*m[2][0]*m[3][1] +
		m[0][3]*m[1][0]*m[2][2]*m[3][1] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
		m[0][2]*m[1][0]*m[2][3]*m[3][1] + m[0][0]*m[1][2]*m[2][3]*m[3][1] +
		m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][1]*m[1][3]*m[2][0]*m[3][2] -
		m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][0]*m[1][3]*m[2][1]*m[3][2] +
		m[0][1]*m[1][0]*m[2][3]*m[3][2] - m[0][0]*m[1][1]*m[2][3]*m[3][2] -
		m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][1]*m[1][2]*m[2][0]*m[3][3] +
		m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][0]*m[1][2]*m[2][1]*m[3][3] -
		m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][0]*m[1][1]*m[2][2]*m[3][3]
}

// robustIncircle finds the determinant of the matrix using the accurate big/Rat type
// |1 x1 y1 x1^2+y1^2|
// |1 x2 y2 x2^2+y2^2|
// |1 x3 y3 x3^2+y3^2|
// |1 x  y  x^2 +y^2 |
func robustIncircle(x1, y1, x2, y2, x3, y3, x, y *big.Rat) bool {
	m := [][]*big.Rat{
		{one, x1, y1, bigAdd(bigMul(x1, x1), bigMul(y1, y1))},
		{one, x2, y2, bigAdd(bigMul(x2, x2), bigMul(y2, y2))},
		{one, x3, y3, bigAdd(bigMul(x3, x3), bigMul(y3, y3))},
		{one, x, y, bigAdd(bigMul(x, x), bigMul(y, y))},
	}
	det := bigAdd(
		bigAdd(
			bigAdd(
				bigAdd(row(3, 2, 1, 0, false, m), row(3, 1, 2, 0, true, m)),
				row(2, 1, 3, 0, false, m),
			),
			bigAdd(
				bigAdd(row(3, 2, 0, 1, true, m), row(3, 0, 2, 1, false, m)),
				row(2, 0, 3, 1, true, m),
			),
		),
		bigAdd(
			bigAdd(
				bigAdd(row(3, 1, 0, 2, false, m), row(3, 0, 1, 2, true, m)),
				row(1, 0, 3, 2, false, m),
			),
			bigAdd(
				bigAdd(row(2, 1, 0, 3, true, m), row(2, 0, 1, 3, false, m)),
				row(1, 0, 2, 3, true, m),
			),
		),
	)
	return det.Cmp(zero) < 0
}

// row is a helper function for robustIncircle
// m[0][a]*m[1][b]*m[2][c]*m[3][d] - m[0][b]*m[1][a]*m[2][c]*m[3][d] /
// - m[0][a]*m[1][b]*m[2][c]*m[3][d] + m[0][b]*m[1][a]*m[2][c]*m[3][d]
func row(a, b, c, d int, plus bool, m [][]*big.Rat) *big.Rat {
	if plus {
		return bigSub(
			bigMul(bigMul(m[0][b], m[1][a]), bigMul(m[2][c], m[3][d])),
			bigMul(bigMul(m[0][a], m[1][b]), bigMul(m[2][c], m[3][d])),
		)
	}
	return bigSub(
		bigMul(bigMul(m[0][a], m[1][b]), bigMul(m[2][c], m[3][d])),
		bigMul(bigMul(m[0][b], m[1][a]), bigMul(m[2][c], m[3][d])),
	)
}
