// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// TiledPlot is a regularly spaced set of plots, aranged as tiles.
type TiledPlot struct {
	Plots []*Plot
	Tiles draw.Tiles
	Align bool // whether to align all tiles axes
}

// NewTiledPlot creates a new set of plots aranged as tiles.
// By default, NewTiledPlot will put a 1 vg.Length space between each plot.
func NewTiledPlot(tiles draw.Tiles) *TiledPlot {
	const pad = 1
	for _, v := range []*vg.Length{
		&tiles.PadTop, &tiles.PadBottom, &tiles.PadRight, &tiles.PadLeft,
		&tiles.PadX, &tiles.PadY,
	} {
		if *v == 0 {
			*v = pad
		}
	}

	plot := &TiledPlot{
		Plots: make([]*Plot, tiles.Rows*tiles.Cols),
		Tiles: tiles,
	}

	for i := 0; i < tiles.Rows; i++ {
		for j := 0; j < tiles.Cols; j++ {
			plot.Plots[i*tiles.Cols+j] = New()
		}
	}

	return plot
}

// Plot returns the plot at the i-th column and j-th row in the set of
// tiles.
// (0,0) is at the top-left of the set of tiles.
func (tp *TiledPlot) Plot(i, j int) *Plot {
	return tp.Plots[i*tp.Tiles.Cols+j]
}

// Draw draws the tiled plot to a draw.Canvas.
//
// Each non-nil plot.Plot in the aranged set of tiled plots is drawn
// inside its dedicated sub-canvas, using hplot.Plot.Draw.
func (tp *TiledPlot) Draw(c draw.Canvas) {
	switch {
	case tp.Align:
		ps := make([][]*plot.Plot, tp.Tiles.Rows)
		for row := 0; row < tp.Tiles.Rows; row++ {
			ps[row] = make([]*plot.Plot, tp.Tiles.Cols)
			for col := range ps[row] {
				p := tp.Plots[row*tp.Tiles.Cols+col]
				if p == nil {
					continue
				}
				ps[row][col] = p.Plot
			}
		}
		cs := plot.Align(ps, tp.Tiles, c)
		for i := 0; i < tp.Tiles.Rows; i++ {
			for j := 0; j < tp.Tiles.Cols; j++ {
				p := ps[i][j]
				if p == nil {
					continue
				}
				p.Draw(cs[i][j])
			}
		}

	default:
		for row := 0; row < tp.Tiles.Rows; row++ {
			for col := 0; col < tp.Tiles.Cols; col++ {
				sub := tp.Tiles.At(c, col, row)
				i := row*tp.Tiles.Cols + col
				p := tp.Plots[i]
				if p == nil {
					continue
				}
				p.Draw(sub)
			}
		}
	}
}

// Save saves the plots to an image file.
// The file format is determined by the extension.
//
// Supported extensions are the same ones than hplot.Plot.Save.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (tp *TiledPlot) Save(w, h vg.Length, file string) error {
	return Save(tp, w, h, file)
}

var (
	_ Drawer = (*TiledPlot)(nil)
)
