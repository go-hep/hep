// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// Band implements the plot.Plotter interface, drawing a colored band made of
// two lines.
type Band struct {
	top    plotter.XYs
	bottom plotter.XYs

	// LineStyle is the style of the line contouring the band.
	// Use zero width to disable.
	draw.LineStyle

	// FillColor is the color to fill the area between
	// the top and bottom data points.
	// Use nil to disable the filling.
	FillColor color.Color
}

func NewBand(fill color.Color, top, bottom plotter.XYer) *Band {
	band := &Band{
		top:       make(plotter.XYs, top.Len()),
		bottom:    make(plotter.XYs, bottom.Len()),
		FillColor: fill,
	}
	for i := range band.top {
		x, y := top.XY(i)
		band.top[i].X = x
		band.top[i].Y = y
	}
	for i := range band.bottom {
		x, y := bottom.XY(i)
		band.bottom[i].X = x
		band.bottom[i].Y = y
	}

	return band
}

func (band *Band) Plot(c draw.Canvas, plt *plot.Plot) {
	switch {
	case len(band.top) <= 1:
		return
	case len(band.bottom) <= 1:
		return
	}

	xys := make(plotter.XYs, 0, len(band.top)+len(band.bottom))
	xys = append(xys, band.bottom...)
	for i := range band.top {
		xys = append(xys, band.top[len(band.top)-1-i])
	}

	poly := plotter.Polygon{
		XYs:       []plotter.XYs{xys},
		LineStyle: band.LineStyle,
		Color:     band.FillColor,
	}

	poly.Plot(c, plt)
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger interface.
func (band *Band) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin1, xmax1, ymin1, ymax1 := plotter.XYRange(band.top)
	xmin2, xmax2, ymin2, ymax2 := plotter.XYRange(band.bottom)

	xmin = math.Min(xmin1, xmin2)
	xmax = math.Max(xmax1, xmax2)
	ymin = math.Min(ymin1, ymin2)
	ymax = math.Max(ymax1, ymax2)

	return xmin, xmax, ymin, ymax
}

var (
	_ plot.Plotter    = (*VertLine)(nil)
	_ plot.Plotter    = (*HorizLine)(nil)
	_ plot.Plotter    = (*Band)(nil)
	_ plot.DataRanger = (*Band)(nil)
)
