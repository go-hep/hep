// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/palette/brewer"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// H2D implements the plotter.Plotter interface,
// drawing a 2-dim histogram of the data.
type H2D struct {
	// H is the histogramming data
	H *hbook.H2D

	// Palette is the color palette used to render
	// the heat map. Palette must not be nil or
	// return a zero length []color.Color.
	Palette palette.Palette

	// InfoStyle is the style of infos displayed for
	// the histogram (entries, mean, rms)
	Infos HInfos

	p *plotter.HeatMap
}

// NewH2D returns a new 2-dim histogram from a hbook.H2D.
func NewH2D(h *hbook.H2D, p palette.Palette) *H2D {
	if p == nil {
		p, _ = brewer.GetPalette(brewer.TypeAny, "RdYlBu", 11)
	}
	return &H2D{
		H:       h,
		Palette: p,
	}
}

func (h *H2D) pltr() *plotter.HeatMap {
	if h.p == nil {
		h.p = plotter.NewHeatMap(h.H.GridXYZ(), h.Palette)
	}
	return h.p
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (h *H2D) Plot(c draw.Canvas, p *plot.Plot) {
	h.pltr().Plot(c, p)
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (h *H2D) DataRange() (xmin, xmax, ymin, ymax float64) {
	return h.pltr().DataRange()
}

// GlyphBoxes returns a slice of GlyphBoxes,
// one for each of the bins, implementing the
// plot.GlyphBoxer interface.
func (h *H2D) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	return h.pltr().GlyphBoxes(p)
}

// check interfaces
var _ plot.Plotter = (*H2D)(nil)
var _ plot.DataRanger = (*H2D)(nil)
var _ plot.GlyphBoxer = (*H2D)(nil)
