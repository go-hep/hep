// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"log"

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
			p := tp.Plot(i, j)
			p.X.Min = -5
			p.X.Max = +5
			newHist(p)
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", i, j)
		}
	}

	// remove plot at (0,1)
	tp.Plots[1] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
