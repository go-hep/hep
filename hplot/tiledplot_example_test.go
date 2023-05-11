// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// An example of making a tile-plot
func ExampleTiledPlot() {
	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 3, Rows: 2})

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	newHist := func(p *hplot.Plot) {
		const npoints = 10000
		hist := hbook.NewH1D(20, -4, +4)
		for i := 0; i < npoints; i++ {
			v := dist.Rand()
			hist.Fill(v, 1)
		}

		h := hplot.NewH1D(hist)
		p.Add(h)
	}

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(j, i)
			p.X.Min = -5
			p.X.Max = +5
			newHist(p)
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", j, i)
		}
	}

	// remove plot at (1,0)
	tp.Plots[1] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}

// An example of making aligned tile-plots
func ExampleTiledPlot_align() {
	tp := hplot.NewTiledPlot(draw.Tiles{
		Cols: 3, Rows: 3,
		PadX: 20, PadY: 20,
	})
	tp.Align = true

	points := func(i, j int) []hbook.Point2D {
		n := i*tp.Tiles.Cols + j + 1
		i += 1
		j = int(math.Pow(10, float64(n)))

		var pts []hbook.Point2D
		for ii := 0; ii < 10; ii++ {
			pts = append(pts, hbook.Point2D{
				X: float64(i + ii),
				Y: float64(j + ii + 1),
			})
		}
		return pts

	}

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(j, i)
			p.X.Min = -5
			p.X.Max = +5
			s := hplot.NewS2D(hbook.NewS2D(points(i, j)...))
			s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
			s.GlyphStyle.Radius = vg.Points(4)
			p.Add(s)

			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", j, i)
		}
	}

	// remove plot at (1,1)
	tp.Plots[4] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_aligned_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
