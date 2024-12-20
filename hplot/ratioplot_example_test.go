// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/vg"
)

func ExampleRatioPlot() {

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewPCG(0, 0)),
	}

	hist1 := hbook.NewH1D(20, -4, +4)
	hist2 := hbook.NewH1D(20, -4, +4)

	for range npoints {
		v1 := dist.Rand() - 0.5
		v2 := dist.Rand() + 0.5
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	rp := hplot.NewRatioPlot()
	rp.Ratio = 0.3

	// Make a plot and set its title.
	rp.Top.Title.Text = "Histos"
	rp.Top.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1 := hplot.NewH1D(hist1)
	h1.FillColor = color.NRGBA{R: 255, A: 100}
	rp.Top.Add(h1)

	h2 := hplot.NewH1D(hist2)
	h2.FillColor = color.NRGBA{B: 255, A: 100}
	rp.Top.Add(h2)

	rp.Top.Add(hplot.NewGrid())

	hist3 := hbook.NewH1D(20, -4, +4)
	for i := range hist3.Len() {
		v1 := hist1.Value(i)
		v2 := hist2.Value(i)
		x1, _ := hist1.XY(i)
		hist3.Fill(x1, v1-v2)
	}

	hdiff := hplot.NewH1D(hist3)

	rp.Bottom.X.Label.Text = "X"
	rp.Bottom.Y.Label.Text = "Delta-Y"
	rp.Bottom.Add(hdiff)
	rp.Bottom.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	err := hplot.Save(rp, width, height, "testdata/diff_plot.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
