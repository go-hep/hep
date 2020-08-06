// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"log"

	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"go-hep.org/x/hep/hplot"
)

func ExampleLabel() {

	// Creating a new plot
	p := hplot.New()
	p.Title.Text = "Plot labels"
	p.X.Min = -1
	p.X.Max = 1
	p.Y.Min = -1
	p.Y.Max = 1

	// Default labels
	l1 := hplot.NewLabel(-0.8, 0.5, "Default label.")
	p.Add(l1)

	// Label with normalized coordinates.
	l3 := hplot.NewLabel(0.5, 0.5, "Label with relative coordinates.",
		hplot.WithNormalized(true),
	)
	p.Add(l3)

	// Label with normalized coordinates and auto-adjustement.
	l4 := hplot.NewLabel(0.95, 0.95, "Label at the canvas edge, with AutoAdjust",
		hplot.WithNormalized(true),
		//hplot.WithAutoAdjust(true),
	)
	p.Add(l4)

	// Label with a customed TextStyle
	usrFont, err := vg.MakeFont("Courier-Bold", 12)
	if err != nil {
		panic(fmt.Errorf("could not create font (Courier-Bold, 12): %w", err))
	}
	sty := draw.TextStyle{
		Color: plotutil.Color(2),
		Font:  usrFont,
	}
	l5 := hplot.NewLabel(0.0, 0.1, "Label with a user-defined font",
		hplot.WithTextStyle(sty),
		hplot.WithNormalized(true),
	)
	p.Add(l5)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/label_plot.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
