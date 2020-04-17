// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type RatioPlot struct {
	Top    *Plot
	Bottom *Plot

	// Tiles controls the layout of the 2x1 ratio-plot grid.
	// Tiles can be used to customize the padding between plots.
	Tiles draw.Tiles

	// Ratio controls how the vertical space is partioned between
	// the top and bottom plots.
	// The top plot will take (1-ratio)*height.
	// Default is 0.3.
	Ratio float64
}

func NewRatioPlot() *RatioPlot {
	rp := &RatioPlot{
		Top:    New(),
		Bottom: New(),
		Ratio:  0.3,
		Tiles:  draw.Tiles{Rows: 2, Cols: 1},
	}

	const pad = 1
	for _, v := range []*vg.Length{
		&rp.Tiles.PadTop, &rp.Tiles.PadBottom,
		&rp.Tiles.PadRight, &rp.Tiles.PadLeft,
		&rp.Tiles.PadX, &rp.Tiles.PadY,
	} {
		if *v == 0 {
			*v = pad
		}
	}

	// hide X-axis labels
	rp.Top.X.Tick.Marker = NoTicks{}
	return rp
}

// Draw draws a ratio plot to a draw.Canvas.
//
// Plotters are drawn in the order in which they were
// added to the plot.  Plotters that  implement the
// GlyphBoxer interface will have their GlyphBoxes
// taken into account when padding the plot so that
// none of their glyphs are clipped.
func (rp *RatioPlot) Draw(dc draw.Canvas) {
	var (
		top, bot = rp.align(dc)
	)

	rp.Top.Draw(top)
	rp.Bottom.Draw(bot)
}

func (rp *RatioPlot) align(dc draw.Canvas) (top, bot draw.Canvas) {
	var (
		ratio = vg.Length(rp.Ratio)
		h     = dc.Size().Y
		ps    = [][]*plot.Plot{
			{rp.Top.Plot},
			{rp.Bottom.Plot},
		}
		cs = plot.Align(ps, rp.Tiles, dc)
	)

	top = cs[0][0]
	bot = cs[1][0]

	top.Rectangle.Min.Y = ratio * h
	top.Rectangle.Max.Y = h
	bot.Rectangle.Max.Y = ratio * h

	return top, bot
}

var (
	_ Drawer = (*RatioPlot)(nil)
)
