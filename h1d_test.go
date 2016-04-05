// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
)

// An example of making a 1D-histogram.
func ExampleHistogram() {
	const npoints = 10000
	var hmax = 1.0

	// stdNorm returns the probability of drawing a
	// value from a standard normal distribution.
	stdNorm := func(x float64) float64 {
		const sigma = 1.0
		const mu = 0.0
		const root2π = 2.50662827459517818309
		return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma)) * hmax
	}
	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := rand.NormFloat64()
		hist.Fill(v, 1)
	}

	// Make a plot and set its title.
	p, err := hplot.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := hplot.NewH1D(hist)
	if err != nil {
		panic(err)
	}
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	// normalize histo
	hmax = h.Hist.Max() / stdNorm(0)

	// The normal distribution function
	norm := hplot.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_plot.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

func TestHistogram1D(t *testing.T) {
	ExampleHistogram()
}
