// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/plot/plotter"
)

func TestH2D(t *testing.T) {
	const (
		nx   = 10
		xmin = 0.0
		xmax = 100.0
		ny   = 10
		ymin = 0.0
		ymax = 100.0
	)

	h := hbook.NewH2D(nx, xmin, xmax, ny, ymin, ymax)
	if h == nil {
		t.Fatalf("nil pointer to H2D")
	}

	if min := h.XMin(); min != xmin {
		t.Errorf("x-min error: got=%v. want=%v\n", min, xmin)
	}
	if max := h.XMax(); max != xmax {
		t.Errorf("x-max error: got=%v. want=%v\n", max, xmax)
	}
	if min := h.YMin(); min != ymin {
		t.Errorf("y-min error: got=%v. want=%v\n", min, ymin)
	}
	if max := h.YMax(); max != ymax {
		t.Errorf("y-max error: got=%v. want=%v\n", max, ymax)
	}

	if name := h.Name(); name != "" {
		t.Errorf("name error: got=%q. want=%q\n", name, "")
	}
	h.Annotation()["name"] = "h1"
	if name := h.Name(); name != "h1" {
		t.Errorf("name error: got=%q. want=%q\n", name, "h1")
	}

	if n := h.Entries(); n != 0 {
		t.Errorf("entries error: got=%v. want=%v\n", n, 0)
	}

	h.Fill(1, 1, 1)
	if n, want := h.Entries(), int64(1); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h.EffEntries(), 1.0; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}

	if w, want := h.SumW(), 1.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h.SumW2(), 1.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h.XMean(), 1.0; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XVariance(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XStdDev(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YMean(), 1.0; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YVariance(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YStdDev(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h.Fill(23, 1, 1)
	if n, want := h.Entries(), int64(2); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h.EffEntries(), 2.0; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}
	if w, want := h.SumW(), 2.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h.SumW2(), 2.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h.XMean(), 12.0; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XVariance(), 242.0; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XStdDev(), 15.556349186104045; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YMean(), 1.0; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YVariance(), 0.0; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YStdDev(), 0.0; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h.Fill(200, 200, 1)
	if w, want := h.SumW(), 3.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h.SumW2(), 3.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h.XMean(), 74.66666666666667; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XVariance(), 11902.333333333334; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XStdDev(), 109.09781543795152; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YMean(), 67.33333333333333; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YVariance(), 13200.333333333334; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YStdDev(), 114.89270356873553; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h.Fill(-100, -100, 0.5)
	if n, want := h.Entries(), int64(4); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h.EffEntries(), 3.769230769230769; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}
	if w, want := h.SumW(), 3.5; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h.SumW2(), 3.25; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h.XMean(), 49.714285714285715; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XVariance(), 14342.111111111111; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.XStdDev(), 119.75855339436558; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YMean(), 43.42857142857143; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YVariance(), 14933.666666666666; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h.YStdDev(), 122.20338238635895; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}
}

func TestH2Edges(t *testing.T) {
	h := hbook.NewH2DFromEdges(
		[]float64{+0, +1, +2, +3},
		[]float64{-3, -2, -1, +0},
	)
	if got, want := h.XMin(), +0.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h.YMin(), -3.0; got != want {
		t.Errorf("got ymin=%v. want=%v", got, want)
	}
	if got, want := h.XMax(), +3.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}
	if got, want := h.YMax(), +0.0; got != want {
		t.Errorf("got ymax=%v. want=%v", got, want)
	}
}

func TestH2EdgesWithPanics(t *testing.T) {
	for _, test := range []struct {
		xs []float64
		ys []float64
	}{
		{
			xs: []float64{0},
			ys: []float64{0, 1},
		},
		{
			xs: []float64{0},
			ys: []float64{0},
		},
		{
			xs: []float64{0, 1, 0.5, 2},
			ys: []float64{0, 1, 2},
		},
		{
			xs: []float64{0, 1, 1},
			ys: []float64{0, 1, 2},
		},
		{
			xs: []float64{0, 1, 0, 1},
			ys: []float64{0, 1, 2},
		},
		{
			xs: []float64{0, 1, 2, 2},
			ys: []float64{0, 1, 2},
		},
		{
			xs: []float64{0, 1, 2, 2, 2},
			ys: []float64{0, 1, 2},
		},
	} {
		{
			panicked, _ := panics(func() {
				_ = hbook.NewH2DFromEdges(test.xs, test.ys)
			})
			if !panicked {
				t.Errorf("edges {x=%v, y=%v} should have panicked", test.xs, test.ys)
			}
		}
		{
			panicked, _ := panics(func() {
				_ = hbook.NewH2DFromEdges(test.ys, test.xs)
			})
			if !panicked {
				t.Errorf("edges {y=%v, x=%v} should have panicked", test.xs, test.ys)
			}
		}
	}
}

// check H2D can be plotted
var _ plotter.GridXYZ = ((*hbook.H2D)(nil)).GridXYZ()

func TestH2DWriteYODA(t *testing.T) {
	h := hbook.NewH2D(5, -1, 1, 5, -2, +2)
	h.Fill(+0.5, +1, 1)
	h.Fill(-0.5, +1, 1)
	h.Fill(+0.0, -1, 1)

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile("testdata/h2d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h2d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestH2DReadYODA(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/h2d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h hbook.H2D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h2d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestH2DBin(t *testing.T) {
	h := hbook.NewH2DFromEdges(
		[]float64{+0, +1, +2, +3},
		[]float64{-3, -2, -1, +0},
	)
	if got, want := h.XMin(), +0.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h.YMin(), -3.0; got != want {
		t.Errorf("got ymin=%v. want=%v", got, want)
	}
	if got, want := h.XMax(), +3.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}
	if got, want := h.YMax(), +0.0; got != want {
		t.Errorf("got ymax=%v. want=%v", got, want)
	}

	h.Fill(0, -3, 1)

	h.Fill(0, -2, 1)
	h.Fill(0, -2, 1)

	h.Fill(1, -2, 1)
	h.Fill(1, -2, 1)
	h.Fill(1, -2, 1)

	for _, tc := range []struct {
		x, y float64
		bin  int
	}{
		{0, -3, 1},
		{0, -2, 2},
		{1, -2, 3},
		{-1, -10, -1},
		{0, -10, -1},
	} {
		t.Run(fmt.Sprintf("x,y=(%v,%v)", tc.x, tc.y), func(t *testing.T) {
			bin := h.Bin(tc.x, tc.y)
			if tc.bin < 0 && bin == nil {
				// ok
				return
			}
			if bin == nil {
				t.Fatalf("unexpected nil bin")
			}

			if bin.EffEntries() != float64(tc.bin) {
				t.Fatalf("x=%v,%v got=%v %v, want=%d", tc.x, tc.y, bin.EffEntries(), bin.Entries(), tc.bin)
			}
		})
	}
}

func TestH2DFillN(t *testing.T) {
	h1 := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	h2 := hbook.NewH2D(10, 0, 10, 10, 0, 10)

	xs := []float64{1, 2, 3, 4}
	ys := []float64{1, 2, 3, 4}
	ws := []float64{1, 2, 1, 1}

	for i := range xs {
		h1.Fill(xs[i], ys[i], ws[i])
	}
	h2.FillN(xs, ys, ws)

	if s1, s2 := h1.SumW(), h2.SumW(); s1 != s2 {
		t.Fatalf("invalid sumw: h1=%v, h2=%v", s1, s2)
	}

	for i := range xs {
		h1.Fill(xs[i], ys[i], 1)
	}
	h2.FillN(xs, ys, nil)

	if s1, s2 := h1.SumW(), h2.SumW(); s1 != s2 {
		t.Fatalf("invalid sumw: h1=%v, h2=%v", s1, s2)
	}

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatalf("expected a panic!")
			}
			const want = "hbook: lengths mismatch"
			if got, want := err.(error).Error(), want; got != want {
				t.Fatalf("invalid panic message:\ngot= %q\nwant=%q", got, want)
			}
		}()

		h2.FillN(xs, ys[:len(xs)-2], nil)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatalf("expected a panic!")
			}
			const want = "hbook: lengths mismatch"
			if got, want := err.(error).Error(), want; got != want {
				t.Fatalf("invalid panic message:\ngot= %q\nwant=%q", got, want)
			}
		}()

		h2.FillN(xs, []float64{1}, ws)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatalf("expected a panic!")
			}
			const want = "hbook: lengths mismatch"
			if got, want := err.(error).Error(), want; got != want {
				t.Fatalf("invalid panic message:\ngot= %q\nwant=%q", got, want)
			}
		}()

		h2.FillN(xs, ys, []float64{1})
	}()
}
