// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"io/ioutil"
	"math"
	"math/rand"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hbook/internal/cmpimg"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/palette/brewer"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat/distmv"
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

	h1 := hbook.NewH2D(nx, xmin, xmax, ny, ymin, ymax)
	if h1 == nil {
		t.Fatalf("nil pointer to H2D")
	}

	if min := h1.MinX(); min != xmin {
		t.Errorf("x-min error: got=%v. want=%v\n", min, xmin)
	}
	if max := h1.MaxX(); max != xmax {
		t.Errorf("x-max error: got=%v. want=%v\n", max, xmax)
	}
	if min := h1.MinY(); min != ymin {
		t.Errorf("y-min error: got=%v. want=%v\n", min, ymin)
	}
	if max := h1.MaxY(); max != ymax {
		t.Errorf("y-max error: got=%v. want=%v\n", max, ymax)
	}

	if name := h1.Name(); name != "" {
		t.Errorf("name error: got=%q. want=%q\n", name, "")
	}
	h1.Annotation()["name"] = "h1"
	if name := h1.Name(); name != "h1" {
		t.Errorf("name error: got=%q. want=%q\n", name, "h1")
	}

	if n := h1.Entries(); n != 0 {
		t.Errorf("entries error: got=%v. want=%v\n", n, 0)
	}

	h1.Fill(1, 1, 1)
	if n, want := h1.Entries(), int64(1); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h1.EffEntries(), 1.0; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}

	if w, want := h1.SumW(), 1.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h1.SumW2(), 1.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h1.MeanX(), 1.0; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceX(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevX(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.MeanY(), 1.0; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceY(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevY(), math.NaN(); !math.IsNaN(v) {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h1.Fill(23, 1, 1)
	if n, want := h1.Entries(), int64(2); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h1.EffEntries(), 2.0; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}
	if w, want := h1.SumW(), 2.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h1.SumW2(), 2.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h1.MeanX(), 12.0; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceX(), 242.0; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevX(), 15.556349186104045; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.MeanY(), 1.0; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceY(), 0.0; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevY(), 0.0; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h1.Fill(200, 200, 1)
	if w, want := h1.SumW(), 3.0; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h1.SumW2(), 3.0; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h1.MeanX(), 74.66666666666667; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceX(), 11902.333333333334; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevX(), 109.09781543795152; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.MeanY(), 67.33333333333333; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceY(), 13200.333333333334; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevY(), 114.89270356873553; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}

	h1.Fill(-100, -100, 0.5)
	if n, want := h1.Entries(), int64(4); n != want {
		t.Errorf("entries error: got=%v. want=%v\n", n, want)
	}
	if n, want := h1.EffEntries(), 3.769230769230769; n != want {
		t.Errorf("eff-entries error: got=%v. want=%v\n", n, want)
	}
	if w, want := h1.SumW(), 3.5; w != want {
		t.Errorf("sum-w: got=%v. want=%v\n", w, want)
	}

	if w2, want := h1.SumW2(), 3.25; w2 != want {
		t.Errorf("sum-w2: got=%v. want=%v\n", w2, want)
	}

	if v, want := h1.MeanX(), 49.714285714285715; v != want {
		t.Errorf("x-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceX(), 14342.111111111111; v != want {
		t.Errorf("x-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevX(), 119.75855339436558; v != want {
		t.Errorf("x-std-dev: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.MeanY(), 43.42857142857143; v != want {
		t.Errorf("y-mean: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.VarianceY(), 14933.666666666666; v != want {
		t.Errorf("y-variance: got=%v. want=%v\n", v, want)
	}

	if v, want := h1.StdDevY(), 122.20338238635895; v != want {
		t.Errorf("y-std-dev: got=%v. want=%v\n", v, want)
	}
}

// check H2D can be plotted
var _ plotter.GridXYZ = ((*hbook.H2D)(nil)).GridXYZ()

func TestH2DPlot(t *testing.T) {
	ExampleNewH2D_plot(t)

	want, err := ioutil.ReadFile("testdata/h2d_plot_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile("testdata/h2d_plot.png")
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := cmpimg.Equal("png", got, want); !ok || err != nil {
		t.Fatalf("error: testdata/h2d_plot.png differ with reference file: %v\n", err)
	}
}

func ExampleNewH2D_plot(t *testing.T) {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		t.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		h.Fill(v[0], v[1], 1)
	}

	p, err := plot.New()
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "Hist-2D"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	plt, err := brewer.GetPalette(brewer.TypeAny, "RdYlBu", 11)
	if err != nil {
		t.Fatal(err)
	}
	p.Add(plotter.NewHeatMap(h.GridXYZ(), plt))
	p.Add(plotter.NewGrid())
	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot.png")
	if err != nil {
		t.Fatal(err)
	}
}
