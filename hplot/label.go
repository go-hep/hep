// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package hplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Label struct {
	Text  string         // Text of the label
	X, Y  float64        // Position of the label
	Style draw.TextStyle // Style of the label

	// Compute the position wrt canvas size
	// X, Y are [0, 1].
	Normalize bool
}

// Plot implements the Plotter interface,
// drawing the label on the canvas.
func (l *Label) Plot(c draw.Canvas, plt *plot.Plot) {

	// FIXME[rmadar]: how to properly propagate TextStyle to
	//                the FillString() command? This includes:
	//                color, font, rotation, X/Yalign and TexHandler.
	if l.Normalize {
		da := plt.DataCanvas(c)
		fnt, err := vg.MakeFont(plotter.DefaultFont, vg.Points(12))
		if err != nil {
			panic("couldn't create font.")
		}
		da.FillString(fnt, vg.Point{X: da.X(l.X), Y: da.Y(l.Y)}, l.Text)

	} else {
		pLabels := plotterLabels(l)
		pLabels.Plot(c, plt)
	}
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (l *Label) DataRange() (xmin, xmax, ymin, ymax float64) {
	pLabels := plotterLabels(l)
	return pLabels.DataRange()
}

// GlyphBoxes returns a GlyphBoxe, corresponding
// to the label, implementing the plot.GlyphBoxer interface.
func (l *Label) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	pLabels := plotterLabels(l)
	return pLabels.GlyphBoxes(p)
}

// Internal helper function to get plotter.Labels type.
func plotterLabels(l *Label) *plotter.Labels {

	// Create of the YXlabels.
	xyL := plotter.XYLabels{
		XYs:    []plotter.XY{{X: l.X, Y: l.Y}},
		Labels: []string{l.Text},
	}

	// Create the plotter.Labels
	labels, err := plotter.NewLabels(xyL)
	if err != nil {
		panic("cannot create plotter.Labels")
	}

	// Add the text styles.
	// FIXME(rmadar): this crashes when l.Style is not
	//                assigned.
	// labels.TextStyle = []draw.TextStyle{l.Style}

	// Return the result
	return labels
}

var (
	_ plot.Plotter    = (*Label)(nil)
	_ plot.DataRanger = (*Label)(nil)
)
