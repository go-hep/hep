hplot
====

[![Build Status](https://drone.io/github.com/go-hep/hplot/status.png)](https://drone.io/github.com/go-hep/hplot/latest)

`hplot` is a WIP package relying on `plotinum` to plot histograms,
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

