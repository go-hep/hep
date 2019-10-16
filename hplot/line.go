// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"

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

	if vline.Left != nil {
		c.SetColor(vline.Left)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: x, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if vline.Right != nil {
		c.SetColor(vline.Right)
		rect := vg.Rectangle{
			Min: vg.Point{X: x, Y: ymin},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}

	if vline.Line.Width != 0 {
		c.StrokeLine2(vline.Line, x, ymin, x, ymax)
	}
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

	if hline.Top != nil {
		c.SetColor(hline.Top)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: y},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if hline.Bottom != nil {
		c.SetColor(hline.Bottom)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: xmax, Y: y},
		}
		c.Fill(rect.Path())
	}

	if hline.Line.Width != 0 {
		c.StrokeLine2(hline.Line, xmin, y, xmax, y)
	}
}

var (
	_ plot.Plotter = (*VertLine)(nil)
	_ plot.Plotter = (*HorizLine)(nil)
)
