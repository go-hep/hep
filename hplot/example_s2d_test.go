// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// ExampleS2D draws some scatter points.
func ExampleS2D() {
	const npoints = 1000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	s2d := hbook.NewS2D()

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		s2d.Fill(hbook.Point2D{X: v[0], Y: v[1]})
	}

	p := hplot.New()
	p.Title.Text = "Scatter-2D"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d)
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(2)

	p.Add(s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d.png")
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleS2D_withErrorBars draws some scatter points
// with their error bars.
func ExampleS2D_withErrorBars() {
	pts := []hbook.Point2D{
		{X: 1, Y: 1, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 2, Y: 2, ErrX: hbook.Range{Min: 0.5, Max: 1.5}, ErrY: hbook.Range{Min: 5, Max: 2}},
	}
	s2d := hbook.NewS2D(pts...)

	p := hplot.New()
	p.Title.Text = "Scatter-2D (with error bars)"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d,
		hplot.WithXErrBars(true),
		hplot.WithYErrBars(true),
		hplot.WithGlyphStyle(draw.GlyphStyle{
			Color:  color.RGBA{R: 255, A: 255},
			Radius: vg.Points(4),
			Shape:  draw.CrossGlyph{},
		}),
	)

	p.Add(s)
	p.Legend.Add("s2d", s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d_errbars.png")
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleS2D_withBand draws some scatter points
// with their error bars and a band
func ExampleS2D_withBand() {
	pts := []hbook.Point2D{
		{X: 1, Y: 1, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 2, Y: 2, ErrY: hbook.Range{Min: 5, Max: 2}},
		{X: 3, Y: 3, ErrY: hbook.Range{Min: 2, Max: 2}},
		{X: 4, Y: 4, ErrY: hbook.Range{Min: 1.2, Max: 2}},
	}
	s2d := hbook.NewS2D(pts...)

	p := hplot.New()
	p.Title.Text = "Scatter-2D (with band)"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d, hplot.WithBand(true), hplot.WithYErrBars(true))
	s.GlyphStyle.Color = color.Black
	s.GlyphStyle.Radius = vg.Points(4)
	s.LineStyle.Width = 1
	s.LineStyle.Dashes = plotutil.Dashes(2)

	p.Add(s)
	p.Legend.Add("s2d", s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d_band.png")
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleS2D_withStepsKind draws some scatter points
// with their error bars, using a step-like style
func ExampleS2D_withStepsKind() {
	pts := []hbook.Point2D{
		{X: 1, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 1, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 2, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 2, ErrY: hbook.Range{Min: 5, Max: 2}},
		{X: 3, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 3, ErrY: hbook.Range{Min: 2, Max: 2}},
		{X: 4, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 4, ErrY: hbook.Range{Min: 1.2, Max: 2}},
	}
	s2d := hbook.NewS2D(pts...)

	p := hplot.New()
	p.Title.Text = "Scatter-2D (with steps)"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d, hplot.WithStepsKind(hplot.HiSteps), hplot.WithYErrBars(true))
	s.GlyphStyle.Color = color.Black
	s.GlyphStyle.Radius = vg.Points(4)
	s.LineStyle.Width = 1
	s.LineStyle.Dashes = plotutil.Dashes(2)

	p.Add(s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d_steps.png")
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleS2D_withSteps_withBand draws some scatter points
// with their error bars, using a step-like style together with a band
func ExampleS2D_withStepsKind_withBand() {
	pts := []hbook.Point2D{
		{X: 1, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 1, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 2, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 5, ErrY: hbook.Range{Min: 5, Max: 2}},
		{X: 3, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 10, ErrY: hbook.Range{Min: 2, Max: 2}},
		{X: 4, ErrX: hbook.Range{Min: 0.5, Max: 0.5}, Y: 15, ErrY: hbook.Range{Min: 1.2, Max: 2}},
	}
	s2d := hbook.NewS2D(pts...)

	p := hplot.New()
	p.Title.Text = "Scatter-2D (with steps and band)"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d, hplot.WithStepsKind(hplot.HiSteps), hplot.WithYErrBars(true), hplot.WithBand(true))
	s.GlyphStyle.Color = color.Black
	s.GlyphStyle.Radius = vg.Points(4)
	s.LineStyle.Width = 1
	s.LineStyle.Dashes = plotutil.Dashes(2)

	p.Add(s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d_steps_band.png")
	if err != nil {
		log.Fatal(err)
	}
}
