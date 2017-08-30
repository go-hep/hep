// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package hplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
)

// GnuplotStyle implements a plot style not much different from the Gnuplot-one.
type GnuplotStyle struct{}

func (s GnuplotStyle) DrawPlot(p *plot.Plot, c draw.Canvas) {
	if p.BackgroundColor != nil {
		c.SetColor(p.BackgroundColor)
		c.Fill(c.Rectangle.Path())
	}
	if p.Title.Text != "" {
		cx := p.DataCanvas(c)
		c.FillText(p.Title.TextStyle, cx.Center().X, c.Max.Y, -0.5, -1, p.Title.Text)
		c.Max.Y -= p.Title.Height(p.Title.Text) - p.Title.Font.Extents().Descent
		c.Max.Y -= p.Title.Padding
	}

	p.X.SanitizeRange()
	x := plot.HorizontalAxis{p.X}
	p.Y.SanitizeRange()
	y := plot.VerticalAxis{p.Y}

	ywidth := y.Size()
	xheight := x.Size()

	xda := plot.PadX(p, draw.Crop(c, ywidth-y.Width-y.Padding, 0, 0, 0))
	yda := plot.PadY(p, draw.Crop(c, 0, xheight-x.Width-x.Padding, 0, 0))

	x.Draw(xda)
	y.Draw(yda)
	xmin := xda.Min.X
	xmax := xda.Max.X
	ymin := yda.Min.Y
	ymax := xda.Max.Y
	xda.StrokeLine2(x.LineStyle, xmin, ymax, xmax, ymax)
	xda.StrokeLine2(x.LineStyle, xmin, ymin, xmax, ymin)
	yda.StrokeLine2(y.LineStyle, xmin, ymin, xmin, ymax)
	yda.StrokeLine2(y.LineStyle, xmax, ymin, xmax, ymax)

	datac := plot.PadY(p, plot.PadX(p, draw.Crop(c, ywidth, xheight, 0, 0)))
	for _, data := range p.Plotters() {
		data.Plot(datac, p)
	}

	p.Legend.Draw(draw.Crop(draw.Crop(c, ywidth, 0, 0, 0), 0, 0, xheight, 0))
}
