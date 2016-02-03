// +build ignore

package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
)

const NPOINTS = 10000

var HMAX = 1.0

func main() {
	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < NPOINTS; i++ {
		v := rand.NormFloat64()
		hist.Fill(v, 1)
	}

	// Make a plot and set its title.
	p, err := hplot.New()
	if err != nil {
		panic(err)
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
	HMAX = h.Hist.Max() / stdNorm(0)

	// The normal distribution function
	norm := hplot.NewFunction(stdNorm)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "hist.png"); err != nil {
		panic(err)
	}
}

// stdNorm returns the probability of drawing a
// value from a standard normal distribution.
func stdNorm(x float64) float64 {
	const sigma = 1.0
	const mu = 0.0
	const root2π = 2.50662827459517818309
	return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma)) * HMAX
}
