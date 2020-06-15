// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

// An example of making a colored binned error band
// from scratch.
func ExampleBinnedErrBand() {

	// Number bins
	nBins := 10

	// Creation of a slice of hbook.Count.
	counts := make([]hbook.Count, nBins)
	for i, xrange := range newBinning(nBins, 0, 10) {
		counts[i].XRange = xrange
		counts[i].Val = float64(i + 1)
		counts[i].Err.Low = 0.1 * counts[i].Val
		counts[i].Err.High = 0.1 * counts[i].Val
	}

	// Set 5th bin to zero
	counts[4].Val, counts[4].Err.Low, counts[4].Err.High = 0, 0, 0

	// Binned error band
	b := &hplot.BinnedErrBand{Counts: counts}
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

func ExampleBinnedErrBand_fromH1D() {

	// Histogram
	h := hbook.NewH1D(20, -5, 5)
	for i := 0; i < 1000; i++ {
		x, w := gauss.Rand(), 1.0
		if 0 < x && x < 0.5 {
			w = 0.0
		}
		h.Fill(x, w)
	}

	hp := hplot.NewH1D(h)
	hp.LineStyle.Width = 0
	hp.FillColor = color.NRGBA{R: 180, G: 180, B: 180, A: 200}

	// Binned error band from the histo
	b := hplot.NewBinnedErrBand(h)
	b.FillColor = color.NRGBA{B: 180, A: 100}
	b.LineStyle.Color = color.NRGBA{B: 100, A: 200}
	b.LineStyle.Width = 1

	// Create a new plot and add the histo and the band.
	p := hplot.New()
	p.Title.Text = "Binned Error Band from H1D"
	p.X.Label.Text = "Binned X"
	p.Y.Label.Text = "Y"
	p.Add(hp)
	p.Add(b)

	// Save the result
	err := p.Save(10*vg.Centimeter, -1, "testdata/binnederrband_fromh1d.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}

// newBinning returns a slice of Range corresponding to
// an equally spaced binning.
func newBinning(n int, xmin, xmax float64) []hbook.Range {
	res := make([]hbook.Range, n)
	dx := (xmax - xmin) / float64(n)
	for i := 0; i < n; i++ {
		lo := xmin + float64(i)*dx
		hi := lo + dx
		res[i].Min = lo
		res[i].Max = hi
	}
	return res
}

var (
	gauss = distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}
)
