// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestIncircle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
		want                         RelativePosition
	}{
		{1, 1, 3, 1, 2, 2, 2, 1.5, Inside},
		{1, 1, 3, 1, 2, 2, 2, 3, Outside},
		{1, 1, 3, 1, 2, 2, 2, 0, On},
		{1, 2, 3, 2, 2, 3, 2.707106, 2.707106, Inside},
		{1, 2, 3, 2, 2, 3, 2.707107, 2.707107, Outside},
		{0, 10, -10, 0, 0, -10, 6, 8, On},
		{0, 10, -10, 0, 0, -10, 6, 7.99, Inside},
		{0, 10, -10, 0, 0, -10, 6, 8.01, Outside},
	}
	for _, test := range tests {
		got := Incircle(test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y)
		if got != test.want {
			t.Fatalf("Incircle(%v,%v,%v,%v,%v,%v,%v,%v) = %v, want = %v", test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y, got, test.want)
		}
	}
}

func BenchmarkSimpleIncircle(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
	}{
		{1, 1, 3, 1, 2, 2, 2, 1.5},
		{1, 1, 3, 1, 2, 2, 2, 3},
		{1, 1, 3, 1, 2, 2, 2, 0},
		{1, 2, 3, 2, 2, 3, 2.707106, 2.707106},
		{1, 2, 3, 2, 2, 3, 2.707107, 2.707107},
		{0, 10, -10, 0, 0, -10, 6, 8},
		{0, 10, -10, 0, 0, -10, 6, 7.99},
		{0, 10, -10, 0, 0, -10, 6, 8.01},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			simpleIncircle(test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y)
		}
	}
}

func BenchmarkMatIncircle(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
	}{
		{1, 1, 3, 1, 2, 2, 2, 1.5},
		{1, 1, 3, 1, 2, 2, 2, 3},
		{1, 1, 3, 1, 2, 2, 2, 0},
		{1, 2, 3, 2, 2, 3, 2.707106, 2.707106},
		{1, 2, 3, 2, 2, 3, 2.707107, 2.707107},
		{0, 10, -10, 0, 0, -10, 6, 8},
		{0, 10, -10, 0, 0, -10, 6, 7.99},
		{0, 10, -10, 0, 0, -10, 6, 8.01},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			matIncircle(test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y)
		}
	}
}

func BenchmarkRobustIncircle(b *testing.B) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
	}{
		{1, 1, 3, 1, 2, 2, 2, 1.5},
		{1, 1, 3, 1, 2, 2, 2, 3},
		{1, 1, 3, 1, 2, 2, 2, 0},
		{1, 2, 3, 2, 2, 3, 2.707106, 2.707106},
		{1, 2, 3, 2, 2, 3, 2.707107, 2.707107},
		{0, 10, -10, 0, 0, -10, 6, 8},
		{0, 10, -10, 0, 0, -10, 6, 7.99},
		{0, 10, -10, 0, 0, -10, 6, 8.01},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			robustIncircle(setBig(test.x1), setBig(test.y1), setBig(test.x2), setBig(test.y2), setBig(test.x3), setBig(test.y3), setBig(test.x), setBig(test.y))
		}
	}
}

// matIncircle computes the determinant of the matrix using the mat package.
// Then it returns the relative position based on the value of the determinant.
// |1 x1 y1 x1^2+y1^2|
// |1 x2 y2 x2^2+y2^2|
// |1 x3 y3 x3^2+y3^2|
// |1 x  y  x^2 +y^2 |
func matIncircle(x1, y1, x2, y2, x3, y3, x, y float64) RelativePosition {
	m := mat.NewDense(4, 4, []float64{1, x1, y1, x1*x1 + y1*y1, 1, x2, y2, x2*x2 + y2*y2, 1, x3, y3, x3*x3 + y3*y3, 1, x, y, x*x + y*y})
	logDet, sign := mat.LogDet(m)
	if math.IsInf(logDet, -1) {
		// logDet is negative infinite and therefore the determinant is 0
		return On
	}
	switch sign {
	case 1:
		return Outside
	case -1:
		return Inside
	}
	return IndeterminatePosition
}
