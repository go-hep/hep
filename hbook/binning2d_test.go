// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "testing"

func TestAxis2DCoords(t *testing.T) {
	bng := newBinning2D(10, -1, 1, 40, -2, +2)
	if nx, want := bng.Nx, 10; nx != want {
		t.Errorf("got nx=%d. want=%d\n", nx, want)
	}
	if ny, want := bng.Ny, 40; ny != want {
		t.Errorf("got ny=%d. want=%d\n", ny, want)
	}
	for i, test := range []struct {
		x, y float64
		want int
	}{
		{x: -1.00, y: -2.00, want: 0},
		{x: +0.01, y: -2.00, want: 5},
		{x: +0.90, y: -2.00, want: 9},
		{x: -1.00, y: -1.90, want: 10},
		{x: +0.01, y: -1.90, want: 15},
		{x: +0.90, y: -1.90, want: 19},
		{x: -1.00, y: -1.00, want: 100},
		{x: +0.01, y: -1.00, want: 105},
		{x: +0.90, y: -1.00, want: 109},
		{x: -0.99, y: +0.01, want: 200},
		{x: +0.01, y: +0.01, want: 205},
		{x: +0.90, y: +0.01, want: 209},
		{x: -0.99, y: +1.01, want: 300},
		{x: +0.01, y: +1.01, want: 305},
		{x: +0.99, y: +1.01, want: 309},
		{x: -0.99, y: +1.9001, want: 390},
		{x: +0.01, y: +1.9001, want: 395},
		{x: +0.99, y: +1.9001, want: bng.Nx*bng.Ny - 1},
		{x: +0.00, y: +2.00, want: -BngN},
		{x: +0.00, y: -2.10, want: -BngS},
		{x: +1.00, y: +2.00, want: -BngNE},
		{x: +1.00, y: -2.00, want: -BngE},
		{x: +1.00, y: -2.10, want: -BngSE},
		{x: -1.10, y: -2.10, want: -BngSW},
		{x: -1.10, y: -2.00, want: -BngW},
		{x: -1.10, y: +2.00, want: -BngNW},
	} {
		got := bng.coordToIndex(test.x, test.y)
		if got != test.want {
			t.Errorf("error: coords[%d](%v, %v). got=%d, want=%d\n", i, test.x, test.y, got, test.want)
			if got >= 0 && test.want >= 0 {
				gbin := bng.Bins[got]
				wbin := bng.Bins[test.want]
				t.Errorf("got.bin.x= %+v\n", gbin.XRange)
				t.Errorf("got.bin.y= %+v\n", gbin.YRange)
				t.Errorf("wnt.bin.x= %+v\n", wbin.XRange)
				t.Errorf("wnt.bin.y= %+v\n", wbin.YRange)
			}
		}
	}
}

func TestAxis2DCoordsFromEdges(t *testing.T) {
	bng := newBinning2DFromEdges(
		[]float64{-1.0, -0.8, -0.6, -0.4, -0.2, 0.0, 0.2, 0.4, 0.6, 0.8, 1.},
		[]float64{
			-2.0, -1.9, -1.8, -1.7, -1.6, -1.5, -1.4, -1.3, -1.2, -1.1,
			-1.0, -0.9, -0.8, -0.7, -0.6, -0.5, -0.4, -0.3, -0.2, -0.1,
			+0.0, +0.1, +0.2, +0.3, +0.4, +0.5, +0.6, +0.7, +0.8, +0.9,
			+1.0, +1.1, +1.2, +1.3, +1.4, +1.5, +1.6, +1.7, +1.8, +1.9,
			+2.0,
		},
	)
	if nx, want := bng.Nx, 10; nx != want {
		t.Errorf("got nx=%d. want=%d\n", nx, want)
	}
	if ny, want := bng.Ny, 40; ny != want {
		t.Errorf("got ny=%d. want=%d\n", ny, want)
	}
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
		{x: +0.9, y: +1.9, want: bng.Nx*bng.Ny - 1},
		{x: +0.0, y: +2.0, want: -BngN},
		{x: +0.0, y: -2.1, want: -BngS},
		{x: +1.0, y: +2.0, want: -BngNE},
		{x: +1.0, y: -2.0, want: -BngE},
		{x: +1.0, y: -2.1, want: -BngSE},
		{x: -1.1, y: -2.1, want: -BngSW},
		{x: -1.1, y: -2.0, want: -BngW},
		{x: -1.1, y: +2.0, want: -BngNW},
	} {
		if got := bng.coordToIndex(test.x, test.y); got != test.want {
			t.Errorf("error: coords[%d](%v, %v). got=%d, want=%d\n", i, test.x, test.y, got, test.want)
		}
	}
}
