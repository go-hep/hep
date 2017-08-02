// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"testing"
)

func TestIncircle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
		want                         RelativePosition
	}{
		{1, 1, 3, 1, 2, 2, 2, 1.5, Inside},
		{1, 1, 3, 1, 2, 2, 2, 3, Outside},
		{1, 2, 3, 2, 2, 3, 2.707106, 2.707106, Inside},
		{1, 2, 3, 2, 2, 3, 2.707107, 2.707107, Outside},
	}
	for _, test := range tests {
		got := Incircle(test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y)
		if got != test.want {
			t.Fatalf("Incircle(%v,%v,%v,%v,%v,%v,%v,%v) = %v, want = %v", test.x1, test.y1, test.x2, test.y2, test.x3, test.y3, test.x, test.y, got, test.want)
		}
	}
}
