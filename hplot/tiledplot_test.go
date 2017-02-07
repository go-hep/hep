// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/stat/distuv"
)

// An example of making a tile-plot
func ExampleTiledPlot(t *testing.T) {
	tp, err := hplot.NewTiledPlot(draw.Tiles{Cols: 3, Rows: 2})
	if err != nil {
		t.Fatalf("error: %v\n", err)
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
				t.Fatalf("error creating histogram (%d,%d): %v\n", i, j, err)
			}
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", i, j)
		}
	}

	// remove plot at (0,1)
	tp.Plots[1] = nil

	err = tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestTiledPlot(t *testing.T) {
	ExampleTiledPlot(t)
	checkPlot(t, "testdata/tiled_plot_histogram_golden.png")
}
