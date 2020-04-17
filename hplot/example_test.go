// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"math"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// An example of a plot + sub-plot
func Example_subplot() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histo
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p1 := hplot.New()
	p1.Title.Text = "Histogram"
	p1.X.Label.Text = "X"
	p1.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	p1.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p1.Add(norm)

	// draw a grid
	p1.Add(hplot.NewGrid())

	// make a second plot which will be diplayed in the upper-right
	// of the previous one
	p2 := hplot.New()
	p2.Title.Text = "Sub plot"
	p2.Add(h)
	p2.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgimg.PngCanvas{Canvas: vgimg.New(width, height)}
	dc := draw.New(c)
	p1.Draw(dc)
	sub := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0.70 * width, Y: 0.50 * height},
			Max: vg.Point{X: 1.00 * width, Y: 1.00 * height},
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

func Example_latexplot() {

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = `Gaussian distribution: $f(x) = \frac{e^{-(x - \mu)^{2}/(2\sigma^{2}) }} {\sigma\sqrt{2\pi}}$`
	p.Y.Label.Text = `$f(x)$`
	p.X.Label.Text = `$x$`

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	h.LineStyle.Color = color.RGBA{R: 255, A: 255}
	h.FillColor = nil
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	p.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	pp := hplot.Wrap(p, hplot.WithBorder(hplot.Border{
		Left:   5,
		Right:  5,
		Top:    5,
		Bottom: 5,
	}))

	err := hplot.Save(pp, width, height, "testdata/latex_plot.tex")
	if err != nil {
		log.Fatalf("could not save LaTeX plot: %+v\n", err)
	}
}

func ExampleSave() {
	p := hplot.New()
	p.Title.Text = "my title"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	const (
		width  = -1 // automatically choose a nice plot width
		height = -1 // automatically choose a nice plot height
	)

	err := hplot.Save(
		p,
		width, height,
		"testdata/plot_save.eps",
		"testdata/plot_save.jpg",
		"testdata/plot_save.pdf",
		"testdata/plot_save.png",
		"testdata/plot_save.svg",
		"testdata/plot_save.tex",
		"testdata/plot_save.tif",
	)

	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}

	// Output:
}
