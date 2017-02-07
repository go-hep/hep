// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "testing"

func TestAxis2DCoords(t *testing.T) {
	ax := newBinning2D(10, -1, 1, 40, -2, +2)
	for i, test := range []struct {
		x, y float64
		want int
	}{
		{x: -1.0, y: -2.0, want: 0},
		{x: +0.0, y: -2.0, want: 5},
		{x: +0.9, y: -2.0, want: 9},
		{x: -1.0, y: -1.9, want: 10},
		{x: +0.0, y: -1.9, want: 15},
		{x: +0.9, y: -1.9, want: 19},
		{x: -1.0, y: -1.0, want: 100},
		{x: +0.0, y: -1.0, want: 105},
		{x: +0.9, y: -1.0, want: 109},
		{x: -1.0, y: +0.0, want: 200},
		{x: +0.0, y: +0.0, want: 205},
		{x: +0.9, y: +0.0, want: 209},
		{x: -1.0, y: +1.0, want: 300},
		{x: +0.0, y: +1.0, want: 305},
		{x: +0.9, y: +1.0, want: 309},
		{x: -1.0, y: +1.9, want: 390},
		{x: +0.0, y: +1.9, want: 395},
		{x: +0.9, y: +1.9, want: ax.nx*ax.ny - 1},
		{x: +0.0, y: +2.0, want: -bngN},
		{x: +0.0, y: -2.1, want: -bngS},
		{x: +1.0, y: +2.0, want: -bngNE},
		{x: +1.0, y: -2.0, want: -bngE},
		{x: +1.0, y: -2.1, want: -bngSE},
		{x: -1.1, y: -2.1, want: -bngSW},
		{x: -1.1, y: -2.0, want: -bngW},
		{x: -1.1, y: +2.0, want: -bngNW},
	} {
		got := ax.coordToIndex(test.x, test.y)
		if got != test.want {
			t.Errorf("error: coords[%d](%v, %v). got=%d, want=%d\n", i, test.x, test.y, got, test.want)
		}
	}
}
