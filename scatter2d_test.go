// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/go-hep/hplot/internal/cmpimg"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat/distmv"
)

// ExampleScatter2D draws some scatter points.
func ExampleScatter2D(t *testing.T) {
	const npoints = 1000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		t.Fatalf("error creating distmv.Normal")
	}

	s2d := hbook.NewScatter2D()

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		s2d.Fill(hbook.Point2D{X: v[0], Y: v[1]})
	}

	p, err := hplot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Scatter-2D"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewScatter2D(s2d)
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(2)

	p.Add(s)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/scatter2d.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestScatter2D(t *testing.T) {
	ExampleScatter2D(t)

	want, err := ioutil.ReadFile("testdata/scatter2d_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile("testdata/scatter2d.png")
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := cmpimg.Equal("png", got, want); !ok || err != nil {
		t.Fatalf("error: testdata/scatter2d.png differ with reference file: %v\n", err)
	}

}

// ExampleScatter2D_withErrorBars draws some scatter points
// with their error bars.
func ExampleScatter2D_withErrorBars(t *testing.T) {
	pts := []hbook.Point2D{
		{X: 1, Y: 1, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 2, Y: 2, ErrX: hbook.Range{Min: 0.5, Max: 1.5}, ErrY: hbook.Range{Min: 5, Max: 2}},
	}
	s2d := hbook.NewScatter2D(pts...)

	p, err := hplot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Scatter-2D (with error bars)"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewScatter2D(s2d)
	s.XErrBars()
	s.YErrBars()
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(4)

	p.Add(s)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/scatter2d_errbars.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestScatter2DWithErrorBars(t *testing.T) {
	ExampleScatter2D_withErrorBars(t)

	want, err := ioutil.ReadFile("testdata/scatter2d_errbars_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile("testdata/scatter2d_errbars.png")
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := cmpimg.Equal("png", got, want); !ok || err != nil {
		t.Fatalf("error: testdata/scatter2d_errbars.png differ with reference file: %v\n", err)
	}

}
