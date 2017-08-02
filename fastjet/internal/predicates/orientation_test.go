// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"testing"
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
		{2.1, 2.1, 1.1, 1.1, 0.1, 0.1, IndeterminateO, Colinear},
		{2.1, 2.1, 1.1, 1.1, 100.1, 100.1, IndeterminateO, Colinear},
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
		if o != test.simple {
			t.Errorf("x1 = %v, y1 = %v, x2 = %v, y2 = %v, x = %v, y = %v, want.Mat = %v. got= %v\n", test.x1, test.y1, test.x2, test.y2, test.x, test.y, test.simple, o)
		}
	}
}
