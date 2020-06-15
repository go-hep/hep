// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hplot"
)

// An example of making a colored band plot
func ExampleBinnedErrBand() {

	// Binning
	bins := newBinning(10, 0, 10)

	// Y values and errors
	vals := make(plotter.Values, 10)
	errs := make(plotter.YErrors, 10)
	for i := range bins {
		vals[i] = float64(i + 1)
		errs[i].Low = 0.1 * vals[i]
		errs[i].High = 0.1 * vals[i]
	}

	// Set 5th bin to zero
	vals[4], errs[4].Low, errs[4].High = 0, 0, 0

	// Binned error band
	b := hplot.NewBinnedErrBand(bins, vals, errs)
	b.FillColor = color.NRGBA{B: 180, A: 200}
	b.LineStyle.Color = color.NRGBA{R: 180, A: 200}
	b.LineStyle.Width = 2

	// Create a new plot and add b
	p := hplot.New()
	p.Title.Text = "Binned Error Band"
	p.X.Label.Text = "Binned X"
	p.Y.Label.Text = "Y"
	p.Add(b)

	// Save the result
	err := p.Save(10*vg.Centimeter, -1, "testdata/binnederrband.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}

// FIX-ME[rmadar]: once the proper type for b.Bins
//                 is selected, one can use the proper
//                 functions.
func newBinning(n int, xmin, xmax float64) [][2]float64 {
	res := make([][2]float64, n)
	dx := (xmax - xmin) / float64(n)
	for i := 0; i < n; i++ {
		lo := xmin + float64(i)*dx
		hi := lo + dx
		res[i] = [2]float64{lo, hi}
	}
	return res
}
