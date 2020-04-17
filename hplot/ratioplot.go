// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"io"
	"math"

	"go-hep.org/x/hep/hplot/htex"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type RatioPlot struct {
	Top    *Plot
	Bottom *Plot

	// Ratio controls how the vertical space is partioned between
	// the top and bottom plots.
	// The top plot will take (1-ratio)*height.
	// Default is 0.3.
	Ratio float64

	// Latex handles the generation of PDFs from .tex files.
	// The default is to use htex.NoopHandler (a no-op).
	// To enable the automatic generation of PDFs, use DefaultHandler:
	//  p := hplot.New()
	//  p.Latex = htex.DefaultHandler
	Latex htex.Handler
}

func NewRatioPlot() *RatioPlot {
	rp := &RatioPlot{
		Top:    New(),
		Bottom: New(),
		Ratio:  0.3,
		Latex:  htex.NoopHandler{},
	}
	// hide X-axis labels
	rp.Top.X.Tick.Marker = NoTicks{}

	return rp
}

func (rp *RatioPlot) LatexHandler() htex.Handler {
	return rp.Latex
}

// Draw draws a ratio plot to a draw.Canvas.
//
// Plotters are drawn in the order in which they were
// added to the plot.  Plotters that  implement the
// GlyphBoxer interface will have their GlyphBoxes
// taken into account when padding the plot so that
// none of their glyphs are clipped.
func (rp *RatioPlot) Draw(dc draw.Canvas) {
	vgtexBorder(dc)

	var (
		ratio  = vg.Length(rp.Ratio)
		width  = dc.Rectangle.Size().X
		height = dc.Rectangle.Size().Y
	)

	top := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: ratio * height},
			Max: vg.Point{X: width, Y: height},
		},
	}
	rp.Top.Draw(top)

	bottom := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: width, Y: ratio * height},
		},
	}
	rp.Bottom.Draw(bottom)
}

// Save saves the plots to an image file.
// The file format is determined by the extension.
//
// Supported extensions are the same ones than hplot.Plot.Save.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (rp *RatioPlot) Save(w, h vg.Length, file string) error {
	return Save(rp, w, h, file)
}

// WriterTo returns an io.WriterTo that will write the plots as
// the specified image format.
//
// Supported formats are the same ones than hplot.Plot.WriterTo
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (rp *RatioPlot) WriterTo(w, h vg.Length, format string) (io.WriterTo, error) {
	switch {
	case w <= 0 && h <= 0:
		w = vgimg.DefaultWidth
		h = vgimg.DefaultWidth / math.Phi
	case w <= 0:
		w = h * math.Phi
	case h <= 0:
		h = w / math.Phi
	}

	c, err := draw.NewFormattedCanvas(w, h, format)
	if err != nil {
		return nil, err
	}
	rp.Draw(draw.New(c))
	return c, nil
}

var (
	_ Drawer       = (*RatioPlot)(nil)
	_ latexHandler = (*RatioPlot)(nil)
)
