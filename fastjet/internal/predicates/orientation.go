// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"fmt"
	"math/big"
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
		panic(fmt.Errorf("predicates: unknown OrientationKind %d", int(o)))
	}
	panic("unreachable")
}

// Orientation returns how the point (x,y) is oriented with respect to
// the line defined by the points (x1,y1) and (x2,y2).
func Orientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	o := simpleOrientation(x1, y1, x2, y2, x, y)
	if o == IndeterminateOrientation {
		// too close to 0 to give a definite answer.
		// Therefore check with more expansive tests.
		o = robustOrientation(setBig(x1), setBig(y1), setBig(x2), setBig(y2), setBig(x), setBig(y))
	}
	return o
}

// simpleOrientation finds the orientation using the float64 type.
// Its accuracy can't be guaranteed, therefore close decisions
// return Indeterminate which signals Orientation that further
// testing is necessary
func simpleOrientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	if (x1 == x2 && x2 == x) || (y1 == y2 && y2 == y) {
		// points are horizontally or vertically aligned
		return Colinear
	}
	// Compute the determinant of the matrix
	//  | x1 y1 1 |
	//  | x2 y2 1 |
	//  | x  y  1 |
	p := newFloat64Pred(x1).mulFloat64(y2).addFloat64Pred(newFloat64Pred(x2).mulFloat64(y)).
		addFloat64Pred(newFloat64Pred(x).mulFloat64(y1)).subFloat64Pred(newFloat64Pred(x1).mulFloat64(y)).
		subFloat64Pred(newFloat64Pred(x2).mulFloat64(y1)).subFloat64Pred(newFloat64Pred(x).mulFloat64(y2))
	// det := x1*y2 + x2*y + x*y1 - x1*y - x2*y1 - x*y2
	det := p.n
	// e determines when the determinant in orientation is too close to 0 to rely on floating point operations.
	// Each intermediate result can have a potential absolute relative rounding error of macheps.
	// If y is the machine representation of x then |(x-y)/x| <= macheps and |x-y| = e, therefore
	// e = macheps*|x|. Since there are no chained multiplications the intermediate results can be add up.
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
func robustOrientation(x1, y1, x2, y2, x, y *big.Rat) OrientationKind {
	// Compute the determinant of the matrix
	//  | x1 y1 1 |
	//  | x2 y2 1 |
	//  | x  y  1 |
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
	default:
		return Colinear
	}
}
