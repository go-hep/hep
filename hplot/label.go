// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Label struct {
	Text      string         // Text of the label
	X, Y      float64        // Position of the label
	TextStyle draw.TextStyle // Text style of the label

	// Normalized indicates whether the label position
	//is in data coordinates  or normalized with regard
	// to the canvas space.
	// When normalized, the label position is assumed
	// to fall in the [0, 1] interval.
	Normalized bool
}

// Plot implements the Plotter interface,
// drawing the label on the canvas.
func (lbl Label) Plot(c draw.Canvas, p *plot.Plot) {
	lbls := lbl.labels(c, p)
	lbls.Plot(c, p)
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (lbl Label) DataRange() (xmin, xmax, ymin, ymax float64) {

	if lbl.Normalized {
		return math.Inf(+1), math.Inf(-1), math.Inf(+1), math.Inf(-1)
	}

	pLabels := lbl.labels(draw.Canvas{}, nil)
	return pLabels.DataRange()
}

// GlyphBoxes returns a GlyphBoxe, corresponding
// to the label, implementing the plot.GlyphBoxer interface.
func (lbl Label) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {

	if lbl.Normalized {
		// FIXME[rmadar]: this crashes and I don't understand why.
		return []plot.GlyphBox{
			{X: lbl.X, Y: lbl.Y, Rectangle: lbl.TextStyle.Rectangle(lbl.Text)},
		}
	}

	return lbl.labels(draw.Canvas{}, p).GlyphBoxes(p)
}

// Internal helper function to get plotter.Labels type.
func (lbl *Label) labels(c draw.Canvas, p *plot.Plot) *plotter.Labels {

	if lbl.TextStyle == (draw.TextStyle{}) {
		// FIXME[rmadar]: implement the hplot.DefaultFont
		//                and hplot.DefaultFontSize.
		defaultFont, err := vg.MakeFont(plotter.DefaultFont, plotter.DefaultFontSize)
		if err != nil {
			panic("impossible to make font.")
		}

		lbl.TextStyle = draw.TextStyle{
			Color: color.Black,
			Font:  defaultFont,
		}
	}

	x := lbl.X
	y := lbl.Y
	if lbl.Normalized {
		dc := p.DataCanvas(c)
		x = float64(dc.X(x))
		y = float64(dc.Y(y))
	}

	xyL := plotter.XYLabels{
		XYs:    []plotter.XY{{X: x, Y: y}},
		Labels: []string{lbl.Text},
	}

	lbls, err := plotter.NewLabels(xyL)
	if err != nil {
		panic("cannot create plotter.Labels")
	}

	lbls.TextStyle = []draw.TextStyle{lbl.TextStyle}

	return lbls
}

var (
	_ plot.Plotter    = (*Label)(nil)
	_ plot.DataRanger = (*Label)(nil)
	_ plot.GlyphBoxer = (*Label)(nil)
)
