// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestOrientation(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x, y float64
		want                 OrientationKind
	}{
		{1, 1, 1, 5, 1, 3, Colinear},
		{2, 2, 10, 2, 7, 2, Colinear},
		{2, 2, 4, 4, 3, 3, Colinear},
		{1, 1, 1, 5, 0, 3, CCW},
		{1, 1, 1, 5, 2, 3, CW},
		{2, 2, 10, 2, 7, 2.001, CCW},
		{2, 2, 10, 2, 7, 1.999, CW},
		{2, 2, 5, 3, 3, 2.334, CCW},
		{2, 2, 5, 3, 3, 2.333, CW},
		{0, 11, 4, 10, 7, 9.251, CCW},
		{0, 11, 4, 10, 7, 9.249, CW},
	}
	for _, test := range tests {
		got := Orientation(test.x1, test.y1, test.x2, test.y2, test.x, test.y)
		if got != test.want {
			t.Fatalf("Orientation(%v,%v,%v,%v,%v,%v) = %v, want = %v", test.x1, test.y1, test.x2, test.y2, test.x, test.y, got, test.want)
		}
	}
}

func TestSimpleVsRobustOrientation(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x, y float64
		simple               OrientationKind
		robust               OrientationKind
	}{
		{2.1, 2.1, 1.1, 1.1, 0.1, 0.1, IndeterminateOrientation, Colinear},
		{2.1, 2.1, 1.1, 1.1, 100.1, 100.1, IndeterminateOrientation, Colinear},
		{2.1, 2.1, 1.1, 1.1, 1000.1, 1000.1, IndeterminateOrientation, Colinear},
		{0.5, 0.5, 12, 12, 24, 24, IndeterminateOrientation, Colinear},
		{1000, 2000, 2000, 3000, 10000, 11000, IndeterminateOrientation, Colinear},
	}
	for _, test := range tests {
		o := simpleOrientation(test.x1, test.y1, test.x2, test.y2, test.x, test.y)
		if o != test.simple {
			t.Errorf("x1 = %v, y1 = %v, x2 = %v, y2 = %v, x = %v, y = %v, want.Simple = %v. got= %v\n", test.x1, test.y1, test.x2, test.y2, test.x, test.y, test.simple, o)
		}
		o = robustOrientation(setBig(test.x1), setBig(test.y1), setBig(test.x2), setBig(test.y2), setBig(test.x), setBig(test.y))
		if o != test.robust {
			t.Errorf("x1 = %v, y1 = %v, x2 = %v, y2 = %v, x = %v, y = %v, want.Robust = %v. got= %v\n", test.x1, test.y1, test.x2, test.y2, test.x, test.y, test.robust, o)
		}
		o = matOrientation(test.x1, test.y1, test.x2, test.y2, test.x, test.y)
		if o != test.robust {
			t.Errorf("x1 = %v, y1 = %v, x2 = %v, y2 = %v, x = %v, y = %v, want.Mat = %v. got= %v\n", test.x1, test.y1, test.x2, test.y2, test.x, test.y, test.robust, o)
		}
	}
}

func BenchmarkSimpleOrientation(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x, y float64
	}{
		{1, 1, 1, 5, 1, 3},
		{2, 2, 10, 2, 7, 2},
		{2, 2, 4, 4, 3, 3},
		{1, 1, 1, 5, 0, 3},
		{1, 1, 1, 5, 2, 3},
		{2, 2, 10, 2, 7, 2.001},
		{2, 2, 10, 2, 7, 1.999},
		{2, 2, 5, 3, 3, 2.334},
		{2, 2, 5, 3, 3, 2.333},
		{0, 11, 4, 10, 7, 9.251},
		{0, 11, 4, 10, 7, 9.249},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			simpleOrientation(test.x1, test.y1, test.x2, test.y2, test.x, test.y)
		}
	}
}

func BenchmarkMatOrientation(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x, y float64
	}{
		{1, 1, 1, 5, 1, 3},
		{2, 2, 10, 2, 7, 2},
		{2, 2, 4, 4, 3, 3},
		{1, 1, 1, 5, 0, 3},
		{1, 1, 1, 5, 2, 3},
		{2, 2, 10, 2, 7, 2.001},
		{2, 2, 10, 2, 7, 1.999},
		{2, 2, 5, 3, 3, 2.334},
		{2, 2, 5, 3, 3, 2.333},
		{0, 11, 4, 10, 7, 9.251},
		{0, 11, 4, 10, 7, 9.249},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			matOrientation(test.x1, test.y1, test.x2, test.y2, test.x, test.y)
		}
	}
}

func BenchmarkRobustOrientation(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x, y float64
	}{
		{1, 1, 1, 5, 1, 3},
		{2, 2, 10, 2, 7, 2},
		{2, 2, 4, 4, 3, 3},
		{1, 1, 1, 5, 0, 3},
		{1, 1, 1, 5, 2, 3},
		{2, 2, 10, 2, 7, 2.001},
		{2, 2, 10, 2, 7, 1.999},
		{2, 2, 5, 3, 3, 2.334},
		{2, 2, 5, 3, 3, 2.333},
		{0, 11, 4, 10, 7, 9.251},
		{0, 11, 4, 10, 7, 9.249},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			robustOrientation(setBig(test.x1), setBig(test.y1), setBig(test.x2), setBig(test.y2), setBig(test.x), setBig(test.y))
		}
	}
}

// matOrientation determines the orientation using the mat package.
//
// It first computes the conditional number of the matrix. When the condition number
// is higher than the Condition Tolerance, then we assume the matrix is singular and
// the determinant is 0. If the determinant is not 0 the sign of the determinant is computed.
//  | x1 y1 1 |
//  | x2 y2 1 |
//  | x  y  1 |
// FIXME once LU.Cond() is exported do the factorization here to improve performance
func matOrientation(x1, y1, x2, y2, x, y float64) OrientationKind {
	if (x1 == x2 && x2 == x) || (y1 == y2 && y2 == y) {
		// points are horizontally or vertically aligned
		return Colinear
	}
	m := mat.NewDense(3, 3, []float64{x1, y1, 1, x2, y2, 1, x, y, 1})
	cond := mat.Cond(m, 1)
	if cond > mat.ConditionTolerance {
		return Colinear
	}
	// Since only the sign is needed LogDet achieves the result in faster time.
	_, sign := mat.LogDet(m)
	switch sign {
	case 1:
		return CCW
	case -1:
		return CW
	}
	return IndeterminateOrientation
}
