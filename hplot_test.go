// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
)

// An example of a plot + sub-plot
func Example_subplot() {
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
	p1, err := hplot.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	p1.Title.Text = "Histogram"
	p1.X.Label.Text = "X"
	p1.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := hplot.NewH1D(hist)
	if err != nil {
		panic(err)
	}
	p1.Add(h)

	// normalize histo
	hmax = h.Hist.Max() / stdNorm(0)

	// The normal distribution function
	norm := hplot.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p1.Add(norm)

	// draw a grid
	p1.Add(hplot.NewGrid())

	// make a second plot which will be diplayed in the upper-right
	// of the previous one
	p2, err := hplot.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	p2.Title.Text = "Sub plot"
	p2.Add(h)
	p2.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgimg.PngCanvas{vgimg.New(width, height)}
	dc := draw.New(c)
	p1.Draw(dc)
	sub := draw.Canvas{
		Canvas: dc,
		Rectangle: draw.Rectangle{
			Min: draw.Point{0.70 * width, 0.50 * height},
			Max: draw.Point{1.00 * width, 1.00 * height},
		},
	}
	p2.Draw(sub)

	f, err := os.Create("testdata/sub_plot.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func TestSubPlot(t *testing.T) {
	Example_subplot()
}
