// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"math"
	"testing"

	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// An example of making a colored band plot
func ExampleBand(t *testing.T) {
	const (
		npoints = 100
		xmax    = 10
	)

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	topData := make(plotter.XYs, npoints)
	botData := make(plotter.XYs, npoints)

	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		x := float64(i+1) / xmax

		v1 := dist.Rand()
		v2 := dist.Rand()

		topData[i].X = x
		topData[i].Y = 1/x + v1 + 10

		botData[i].X = x
		botData[i].Y = math.Log(x) + v2
	}

	top, err := hplot.NewLine(topData)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	top.LineStyle.Color = color.RGBA{R: 255, A: 255}

	bot, err := hplot.NewLine(botData)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	bot.LineStyle.Color = color.RGBA{B: 255, A: 255}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 2})

	tp.Plots[0].Title.Text = "Band"
	tp.Plots[0].Add(
		top,
		bot,
		hplot.NewBand(color.Gray{200}, topData, botData),
	)

	tp.Plots[1].Title.Text = "Band"
	var (
		blue = color.RGBA{B: 255, A: 255}
		grey = color.Gray{200}
		band = hplot.NewBand(grey, topData, botData)
	)
	band.LineStyle = plotter.DefaultLineStyle
	band.LineStyle.Color = blue
	tp.Plots[1].Add(band)

	err = tp.Save(10*vg.Centimeter, -1, "testdata/band.png")
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
}

func TestBand(t *testing.T) {
	ExampleBand(t)
	checkPlot(t, "testdata/band_golden.png")
}
