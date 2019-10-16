// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

// An example of making a vertical-line plot
func ExampleVLine() {
	p := hplot.New()
	p.Title.Text = "vlines"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	var (
		left  = color.RGBA{B: 255, A: 255}
		right = color.RGBA{R: 255, A: 255}
	)

	p.Add(
		hplot.VLine(2.5, left, nil),
		hplot.VLine(5, nil, nil),
		hplot.VLine(7.5, nil, right),
	)

	err := p.Save(10*vg.Centimeter, -1, "testdata/vline.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}

// An example of making a horizontal-line plot
func ExampleHLine() {
	p := hplot.New()
	p.Title.Text = "hlines"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	var (
		top    = color.RGBA{B: 255, A: 255}
		bottom = color.RGBA{R: 255, A: 255}
	)

	p.Add(
		hplot.HLine(2.5, nil, bottom),
		hplot.HLine(5, nil, nil),
		hplot.HLine(7.5, top, nil),
	)

	err := p.Save(10*vg.Centimeter, -1, "testdata/hline.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
