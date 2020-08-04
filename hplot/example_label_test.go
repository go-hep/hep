// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"log"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hplot"
)

func ExampleLabel() {

	// Label definition
	l1 := hplot.Label{
		Text: "You are in ...",
		X:    -0.5, Y: 0.5,
	}
	l2 := hplot.Label{
		Text: "... the right place",
		X:    0.5, Y: -0.5,
	}
	l3 := hplot.Label{
		Text: "This is the middle",
		X:    0.5, Y: 0.5,
		Normalized: true}

	// New plot
	p := hplot.New()
	p.Title.Text = "Plot labels"
	p.X.Min = -1
	p.X.Max = 1
	p.Y.Min = -1
	p.Y.Max = 1

	// Adding the labels
	p.Add(l1)
	p.Add(l2)
	p.Add(l3)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/label_plot.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
