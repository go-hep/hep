// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"log"

	xfnt "golang.org/x/image/font"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hplot"
)

func ExampleLabel() {

	// Creating a new plot
	p := hplot.New()
	p.Title.Text = "Plot labels"
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = -10
	p.Y.Max = +10

	// Default labels
	l1 := hplot.NewLabel(-8, 5, "(-8,5)\nDefault label")
	p.Add(l1)

	// Label with normalized coordinates.
	l3 := hplot.NewLabel(
		0.5, 0.5,
		"(0.5,0.5)\nLabel with relative coords",
		hplot.WithLabelNormalized(true),
	)
	p.Add(l3)

	// Label with normalized coordinates and auto-adjustement.
	l4 := hplot.NewLabel(
		0.95, 0.95,
		"(0.95,0.95)\nLabel at the canvas edge, with AutoAdjust",
		hplot.WithLabelNormalized(true),
		hplot.WithLabelAutoAdjust(true),
	)
	p.Add(l4)

	// Label with a customed TextStyle
	usrFont := font.Font{
		Typeface: "Liberation",
		Variant:  "Mono",
		Weight:   xfnt.WeightBold,
		Style:    xfnt.StyleNormal,
		Size:     12,
	}
	sty := text.Style{
		Color: plotutil.Color(2),
		Font:  usrFont,
	}
	l5 := hplot.NewLabel(
		0.0, 0.1,
		"(0.0,0.1)\nLabel with a user-defined font",
		hplot.WithLabelTextStyle(sty),
		hplot.WithLabelNormalized(true),
	)
	p.Add(l5)

	p.Add(plotter.NewGlyphBoxes())
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	err := p.Save(15*vg.Centimeter, -1, "testdata/label_plot.png")
	if err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
