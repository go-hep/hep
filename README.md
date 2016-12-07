hplot
====

[![Build Status](https://secure.travis-ci.org/go-hep/hplot.png)](http://travis-ci.org/go-hep/hplot) [![GoDoc](https://godoc.org/github.com/go-hep/hplot?status.svg)](https://godoc.org/github.com/go-hep/hplot)

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

![hist-example](https://github.com/go-hep/hplot/raw/master/testdata/h1d_plot_golden.png)

```go
package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat/distuv"
)

func main() {
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
	hist.Scale(1 / hist.Integral())

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

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
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
	"github.com/gonum/stat/distuv"
)

// An example of making a tile-plot
func main() {
	tp, err := hplot.NewTiledPlot(draw.Tiles{Cols: 3, Rows: 2})
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:     0,
		Sigma:  1,
		Source: rand.New(rand.NewSource(0)),
	}

	newHist := func(p *hplot.Plot) error {
		const npoints = 10000
		hist := hbook.NewH1D(20, -4, +4)
		for i := 0; i < npoints; i++ {
			v := dist.Rand()
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

### Subplots

![sub-plot](https://github.com/go-hep/hplot/raw/master/testdata/sub_plot_golden.png)

https://godoc.org/github.com/go-hep/hplot#example-package--Subplot

### Diff-plots

![diff-plot](https://github.com/go-hep/hplot/raw/master/testdata/diff_plot_golden.png)

https://godoc.org/github.com/go-hep/hplot#example-package--Diffplot

### LaTeX-plots

[latex-plot (PDF)](https://github.com/go-hep/hplot/blob/master/testdata/latex_plot_golden.pdf)

https://godoc.org/github.com/go-hep/hplot#example-package--Latexplot

### 2D histogram

```go
package main

import (
	"log"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/palette/brewer"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat/distmv"
)

func main() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
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
		log.Fatalf("error: %v\n", err)
	}
	p.Title.Text = "Hist-2D"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	p.Add(hplot.NewH2D(h, nil))
	p.Add(plotter.NewGrid())
	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot.png")
	if err != nil {
		log.Fatal(err)
	}
}
```

![h2d-example](https://github.com/go-hep/hplot/raw/master/testdata/h2d_plot_golden.png)
