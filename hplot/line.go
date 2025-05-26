// Copyright ©2019 The go-hep Authors. All rights reserved.
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

// VertLine draws a vertical line at X and colors the
// left and right portions of the plot with the provided
// colors.
type VertLine struct {
	X     float64
	Line  draw.LineStyle
	Left  color.Color
	Right color.Color
}

// VLine creates a vertical line at x with the default line style.
func VLine(x float64, left, right color.Color) *VertLine {
	return &VertLine{
		X:     x,
		Line:  plotter.DefaultLineStyle,
		Left:  left,
		Right: right,
	}
}

func (vline *VertLine) Plot(c draw.Canvas, plt *plot.Plot) {
	var (
		trX, _ = plt.Transforms(&c)
		x      = trX(vline.X)
		xmin   = c.Min.X
		xmax   = c.Max.X
		ymin   = c.Min.Y
		ymax   = c.Max.Y
	)

	if vline.Left != nil && x > xmin {
		c.SetColor(vline.Left)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: x, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if vline.Right != nil && x < xmax {
		c.SetColor(vline.Right)
		rect := vg.Rectangle{
			Min: vg.Point{X: x, Y: ymin},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}

	if vline.Line.Width != 0 && xmin <= x && x <= xmax {
		c.StrokeLine2(vline.Line, x, ymin, x, ymax)
	}
}

// Thumbnail returns the thumbnail for the VertLine,
// implementing the plot.Thumbnailer interface.
func (vline *VertLine) Thumbnail(c *draw.Canvas) {
	if vline.Left != nil {
		minX := c.Min.X
		maxX := c.Center().X
		minY := c.Min.Y
		maxY := c.Max.Y
		points := []vg.Point{
			{X: minX, Y: minY},
			{X: minX, Y: maxY},
			{X: maxX, Y: maxY},
			{X: maxX, Y: minY},
		}
		poly := c.ClipPolygonY(points)
		c.FillPolygon(vline.Left, poly)
	}

	if vline.Right != nil {
		minX := c.Center().X
		maxX := c.Max.X
		minY := c.Min.Y
		maxY := c.Max.Y
		points := []vg.Point{
			{X: minX, Y: minY},
			{X: minX, Y: maxY},
			{X: maxX, Y: maxY},
			{X: maxX, Y: minY},
		}
		poly := c.ClipPolygonY(points)
		c.FillPolygon(vline.Right, poly)
	}

	if vline.Line.Width != 0 {
		x := c.Center().X
		c.StrokeLine2(vline.Line, x, c.Min.Y, x, c.Max.Y)
	}
}

// DataRange returns the range of X and Y values.
func (vline *VertLine) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = vline.X
	xmax = vline.X
	ymin = math.Inf(+1)
	ymax = math.Inf(-1)
	return
}

// HorizLine draws a horizontal line at Y and colors the
// top and bottom portions of the plot with the provided
// colors.
type HorizLine struct {
	Y      float64
	Line   draw.LineStyle
	Top    color.Color
	Bottom color.Color
}

// HLine creates a horizontal line at y with the default line style.
func HLine(y float64, top, bottom color.Color) *HorizLine {
	return &HorizLine{
		Y:      y,
		Line:   plotter.DefaultLineStyle,
		Top:    top,
		Bottom: bottom,
	}
}

func (hline *HorizLine) Plot(c draw.Canvas, plt *plot.Plot) {
	var (
		_, trY = plt.Transforms(&c)
		y      = trY(hline.Y)
		xmin   = c.Min.X
		xmax   = c.Max.X
		ymin   = c.Min.Y
		ymax   = c.Max.Y
	)

	if hline.Top != nil && y < ymax {
		c.SetColor(hline.Top)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: y},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if hline.Bottom != nil && y > ymin {
		c.SetColor(hline.Bottom)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: xmax, Y: y},
		}
		c.Fill(rect.Path())
	}

	if hline.Line.Width != 0 && ymin <= y && y <= ymax {
		c.StrokeLine2(hline.Line, xmin, y, xmax, y)
	}
}

// Thumbnail returns the thumbnail for the VertLine,
// implementing the plot.Thumbnailer interface.
func (hline *HorizLine) Thumbnail(c *draw.Canvas) {
	if hline.Top != nil {
		minX := c.Min.X
		maxX := c.Max.X
		minY := c.Center().Y
		maxY := c.Max.Y
		points := []vg.Point{
			{X: minX, Y: minY},
			{X: minX, Y: maxY},
			{X: maxX, Y: maxY},
			{X: maxX, Y: minY},
		}
		poly := c.ClipPolygonY(points)
		c.FillPolygon(hline.Top, poly)
	}

	if hline.Bottom != nil {
		minX := c.Min.X
		maxX := c.Max.X
		minY := c.Min.Y
		maxY := c.Center().Y
		points := []vg.Point{
			{X: minX, Y: minY},
			{X: minX, Y: maxY},
			{X: maxX, Y: maxY},
			{X: maxX, Y: minY},
		}
		poly := c.ClipPolygonY(points)
		c.FillPolygon(hline.Bottom, poly)
	}

	if hline.Line.Width != 0 {
		y := c.Center().Y
		c.StrokeLine2(hline.Line, c.Min.X, y, c.Max.X, y)
	}
}

// DataRange returns the range of X and Y values.
func (hline *HorizLine) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(+1)
	xmax = math.Inf(-1)
	ymin = hline.Y
	ymax = hline.Y
	return
}

var (
	_ plot.Plotter = (*VertLine)(nil)
	_ plot.Plotter = (*HorizLine)(nil)

	_ plot.DataRanger = (*VertLine)(nil)
	_ plot.DataRanger = (*HorizLine)(nil)

	_ plot.Thumbnailer = (*VertLine)(nil)
	_ plot.Thumbnailer = (*HorizLine)(nil)
)
