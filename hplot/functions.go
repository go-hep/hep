// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Function implements the Plotter interface,
// drawing a line for the given function.
type Function struct {
	F func(x float64) (y float64)

	// XMin and XMax specify the range
	// of x values to pass to F.
	XMin, XMax float64

	Samples int

	draw.LineStyle

	// LogY allows rendering with a log-scaled Y axis.
	// When enabled, function values returning 0 will be discarded from
	// the final plot.
	LogY bool
}

// NewFunction returns a Function that plots F using
// the default line style with 50 samples.
func NewFunction(f func(float64) float64) *Function {
	return &Function{
		F:         f,
		Samples:   50,
		LineStyle: plotter.DefaultLineStyle,
	}
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (f *Function) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)

	min, max := f.XMin, f.XMax
	if min == 0 && max == 0 {
		min = p.X.Min
		max = p.X.Max
	}
	d := (max - min) / float64(f.Samples-1)
	switch {
	case f.LogY:
		var (
			line  = 0
			lines = [][]vg.Point{make([]vg.Point, 0, f.Samples)}
		)
		for i := range f.Samples {
			x := min + float64(i)*d
			y := f.F(x)
			switch {
			case math.IsInf(y, -1) || y <= 0:
				line++
				lines = append(lines, make([]vg.Point, 0, f.Samples-i))
			default:
				lines[line] = append(lines[line], vg.Point{
					X: trX(x),
					Y: trY(y),
				})
			}
		}
		for _, line := range lines {
			if len(line) <= 1 {
				// FIXME(sbinet): we should find a couple of points around...
				continue
			}
			c.StrokeLines(f.LineStyle, c.ClipLinesXY(line)...)
		}
	default:
		line := make([]vg.Point, f.Samples)
		for i := range line {
			x := min + float64(i)*d
			y := f.F(x)
			line[i].X = trX(x)
			line[i].Y = trY(y)
		}
		c.StrokeLines(f.LineStyle, c.ClipLinesXY(line)...)
	}
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the function.
func (f Function) Thumbnail(c *draw.Canvas) {
	y := c.Center().Y
	c.StrokeLine2(f.LineStyle, c.Min.X, y, c.Max.X, y)
}
