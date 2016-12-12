// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/go-hep/hplot/internal/cmpimg"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat/distuv"
)

// An example of making a 1D-histogram.
func ExampleH1D(t *testing.T) {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:     0,
		Sigma:  1,
		Source: rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning().Bins() {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p, err := hplot.New()
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := hplot.NewH1D(hist)
	if err != nil {
		t.Fatal(err)
	}
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_plot.png"); err != nil {
		t.Fatalf("error saving plot: %v\n", err)
	}
}

func TestHistogram1D(t *testing.T) {
	ExampleH1D(t)

	want, err := ioutil.ReadFile("testdata/h1d_plot_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile("testdata/h1d_plot.png")
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := cmpimg.Equal("png", got, want); !ok || err != nil {
		t.Fatalf("error: testdata/h1d_plot.png differ with reference file: %v\n", err)
	}
}
