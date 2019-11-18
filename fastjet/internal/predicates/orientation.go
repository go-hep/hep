// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"math/big"

	"golang.org/x/xerrors"
	"gonum.org/v1/gonum/mat"
)

// OrientationKind indicates how three points are located in respect to each other.
type OrientationKind int

const (
	// Counterclockwise
	CCW OrientationKind = iota
	// Clockwise
	CW
	Colinear
	IndeterminateOrientation
)

func (o OrientationKind) String() string {
	switch o {
	case CCW:
		return "Counterclockwise"
	case CW:
		return "Clockwise"
	case Colinear:
		return "Colinear"
	case IndeterminateOrientation:
		return "Indeterminate"
	default:
		panic(xerrors.Errorf("predicates: unknown OrientationKind %d", int(o)))
	}
}

// Orientation returns how the point (x,y) is oriented with respect to
// the line defined by the points (x1,y1) and (x2,y2).
func Orientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	o := simpleOrientation(x1, y1, x2, y2, x, y)
	if o == IndeterminateOrientation {
		// too close to 0 to give a definite answer.
		// Therefore check with more expansive tests.
		o = matOrientation(x1, y1, x2, y2, x, y)
	}
	return o
}

// simpleOrientation finds the orientation using the simple float64 type.
// Its accuracy can't be guaranteed, therefore close decisions
// return IndeterminateOrientation which signals Orientation that further
// testing is necessary.
//
// It computes the determinant of the matrix and returns the orientation based
// on the value of the determinant.
//  | x1 y1 1 |
//  | x2 y2 1 |
//  | x  y  1 |
func simpleOrientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	if (x1 == x2 && x2 == x) || (y1 == y2 && y2 == y) {
		// points are horizontally or vertically aligned
		return Colinear
	}
	// Compute the determinant of the matrix
	p := newFloat64Pred(x1).mulFloat64(y2).addFloat64Pred(newFloat64Pred(x2).mulFloat64(y)).
		addFloat64Pred(newFloat64Pred(x).mulFloat64(y1)).subFloat64Pred(newFloat64Pred(x1).mulFloat64(y)).
		subFloat64Pred(newFloat64Pred(x2).mulFloat64(y1)).subFloat64Pred(newFloat64Pred(x).mulFloat64(y2))
	// det := x1*y2 + x2*y + x*y1 - x1*y - x2*y1 - x*y2
	det := p.n
	// e determines when the determinant in simpleOrientation is too close to 0 to rely on floating point operations.
	e := p.e
	if det > e {
		return CCW
	}
	if det < -e {
		return CW
	}
	return IndeterminateOrientation
}

// robustOrientation finds the orientation using the accurate big/Rat type.
//
// It computes the determinant of the matrix and returns the orientation based
// on the value of the determinant.
//  | x1 y1 1 |
//  | x2 y2 1 |
//  | x  y  1 |
func robustOrientation(x1, y1, x2, y2, x, y *big.Rat) OrientationKind {
	// Compute the determinant of the matrix
	// det := x1*y2 + x2*y + x*y1 - x1*y - x2*y1 - x*y2
	det := bigSub(
		bigAdd(bigAdd(bigMul(x1, y2), bigMul(x2, y)), bigMul(x, y1)),
		bigAdd(bigAdd(bigMul(x1, y), bigMul(x2, y1)), bigMul(x, y2)),
	)
	sign := det.Sign()
	switch sign {
	case 1:
		return CCW
	case -1:
		return CW
	case 0:
		return Colinear
	}
	return IndeterminateOrientation
}

// matOrientation determines the orientation using the mat package.
//
// It first computes the conditional number of the matrix. When the condition number
// is higher than the Condition Tolerance, then we assume the matrix is singular and
// the determinant is 0. If the determinant is not 0 the sign of the determinant is computed.
//  | x1 y1 1 |
//  | x2 y2 1 |
//  | x  y  1 |
func matOrientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	if (x1 == x2 && x2 == x) || (y1 == y2 && y2 == y) {
		// points are horizontally or vertically aligned
		return Colinear
	}
	m := mat.NewDense(3, 3, []float64{x1, y1, 1, x2, y2, 1, x, y, 1})
	var lu mat.LU
	lu.Factorize(m)
	cond := lu.Cond()
	if cond > mat.ConditionTolerance {
		return Colinear
	}
	// Since only the sign is needed LogDet achieves the result in faster time.
	_, sign := lu.LogDet()
	switch sign {
	case 1:
		return CCW
	case -1:
		return CW
	}
	return IndeterminateOrientation
}
