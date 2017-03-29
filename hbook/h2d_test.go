// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/stat/distmv"
	"go-hep.org/x/hep/hbook"
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

func ExampleH2D() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		h.Fill(v[0], v[1], 1)
	}
}

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
