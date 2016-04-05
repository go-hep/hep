hplot
====

[![Build Status](https://drone.io/github.com/go-hep/hplot/status.png)](https://drone.io/github.com/go-hep/hplot/latest)

`hplot` is a WIP package relying on `gonum/plot` to plot histograms,
n-tuples and functions.

## Installation

```sh
$ go get github.com/go-hep/hplot
```

## Documentation

Is available on ``godoc``:

http://godoc.org/github.com/go-hep/hplot


## Examples

### 1D histogram

![hist-example](https://github.com/go-hep/hplot/raw/master/examples/hist.png)

```go
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/go-hep/hplot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

const NPOINTS = 10000

var HMAX = 1.0

func main() {
	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))
	v := make(plotter.Values, NPOINTS)
	for i := range v {
		v[i] = rand.NormFloat64()
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
	h, err := hplot.NewHist(v, 16)
	if err != nil {
		panic(err)
	}
	// h.Infos.Style = hplot.HInfo_None
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
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "hist.png"); err != nil {
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
```

### Tiles of 1D histograms

![tiled-plot](https://github.com/go-hep/hplot/raw/master/testdata/tiled_plot_histogram_golden.png)

```go
package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

func main() {
	tp, err := hplot.NewTiledPlot(draw.Tiles{Cols: 3, Rows: 2})
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// Draw some random values from the standard
	// normal distribution.
	rand.Seed(int64(0))

	newHist := func(p *hplot.Plot) error {
		const npoints = 10000
		hist := hbook.NewH1D(20, -4, +4)
		for i := 0; i < npoints; i++ {
			v := rand.NormFloat64()
			hist.Fill(v, 1)
		}

		h, err := hplot.NewH1D(hist)
		if err != nil {
			return err
		}
		p.Add(h)
		return nil
	}

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(i, j)
			p.X.Min = -5
			p.X.Max = +5
			err := newHist(p)
			if err != nil {
				log.Fatalf("error creating histogram (%d,%d): %v\n", i, j, err)
			}
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", i, j)
		}
	}

	// remove plot at (0,1)
	tp.Plots[1] = nil

	err = tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
```
