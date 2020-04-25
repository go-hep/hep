// Copyright ©2016 The go-hep Authors. All rights reserved.
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

// S2D plots a set of 2-dim points with error bars.
type S2D struct {
	Data plotter.XYer

	// GlyphStyle is the style of the glyphs drawn
	// at each point.
	draw.GlyphStyle

	// LineStyle is the style of the line drawn
	// connecting each point.
	// Use zero width to disable.
	LineStyle draw.LineStyle

	xbars *plotter.XErrorBars
	ybars *plotter.YErrorBars

	// Band displays a colored band between the y-min and y-max error bars.
	Band *Band

	// Step enable a step-like plotting style
	Steps StepsKind
}

// withXErrBars enables the X error bars
func (pts *S2D) withXErrBars() error {
	xerr, ok := pts.Data.(plotter.XErrorer)
	if !ok {
		return nil
	}

	type xerrT struct {
		plotter.XYer
		plotter.XErrorer
	}
	xplt, err := plotter.NewXErrorBars(xerrT{pts.Data, xerr})
	if err != nil {
		return err
	}

	pts.xbars = xplt
	return nil
}

// withYErrBars enables the Y error bars
func (pts *S2D) withYErrBars() error {
	yerr, ok := pts.Data.(plotter.YErrorer)
	if !ok {
		return nil
	}

	type yerrT struct {
		plotter.XYer
		plotter.YErrorer
	}
	yplt, err := plotter.NewYErrorBars(yerrT{pts.Data, yerr})
	if err != nil {
		return err
	}

	pts.ybars = yplt
	return nil
}

// withBand enables the band between ymin-ymax error bars.
func (pts *S2D) withBand() error {
	yerr, ok := pts.Data.(plotter.YErrorer)
	if !ok {
		return nil
	}

	var (
		top = make(plotter.XYs, pts.Data.Len())
		bot = make(plotter.XYs, pts.Data.Len())
	)

	switch pts.Steps {
	case NoSteps:
		for i := range top {
			x, y := pts.Data.XY(i)
			ymin, ymax := yerr.YError(i)
			top[i].X = x
			top[i].Y = y + math.Abs(ymax)
			bot[i].X = x
			bot[i].Y = y - math.Abs(ymin)
		}
	case HiSteps:
		for i := range top {
			// WIP(rmadar): implement the proper band in case of step
			//              this might actually involve 2xn points - need to think
			x, y := pts.Data.XY(i)
			ymin, ymax := yerr.YError(i)
			top[i].X = x
			top[i].Y = y + math.Abs(ymax)
			bot[i].X = x
			bot[i].Y = y - math.Abs(ymin)
		}
	}
	pts.Band = NewBand(color.Gray{200}, top, bot)
	return nil
}

// NewS2D creates a 2-dim scatter plot from a XYer.
func NewS2D(data plotter.XYer, opts ...Options) *S2D {
	s := &S2D{
		Data:       data,
		GlyphStyle: plotter.DefaultGlyphStyle,
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}

	cfg := newConfig(opts)

	// rmadar: not sure about this (ie, best to handle default value)
	s.Steps = cfg.steps
	
	if cfg.bars.xerrs {
		_ = s.withXErrBars()
	}

	if cfg.bars.yerrs {
		_ = s.withYErrBars()
	}

	if cfg.band {
		_ = s.withBand()
	}

	if cfg.glyph != (draw.GlyphStyle{}) {
		s.GlyphStyle = cfg.glyph
	}

	return s
}

// Plot draws the Scatter, implementing the plot.Plotter
// interface.
func (pts *S2D) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	if pts.Band != nil {
		pts.Band.Plot(c, plt)
	}

	for i := 0; i < pts.Data.Len(); i++ {
		x, y := pts.Data.XY(i)
		c.DrawGlyph(pts.GlyphStyle, vg.Point{X: trX(x), Y: trY(y)})
	}

	if pts.LineStyle.Width > 0 {
		
		data, err := plotter.CopyXYs(pts.Data)
		if err != nil {
			panic(err)
		}
		
		// rmadar: a switch was suggested but I'd say a if seems more suitable
		if pts.Steps == HiSteps && pts.xbars != nil {
			
			xerr, ok := pts.Data.(plotter.XErrorer)
			if !ok {
				panic("s2d: cannot get X errors during HiSteps plotting")
			}
			
			data_step := plotter.XYs{}
			for i, d := range data {
				xmin, xmax := xerr.XError(i)
				data_step = append(data_step, plotter.XY{X: d.X - xmin, Y: d.Y} )
				data_step = append(data_step, plotter.XY{X: d.X + xmax, Y: d.Y} )
			}
		}
		
		line := plotter.Line{
			XYs:       data,
			LineStyle: pts.LineStyle,
		}
		line.Plot(c, plt)
	}

	if pts.xbars != nil {
		pts.xbars.LineStyle.Color = pts.GlyphStyle.Color
		pts.xbars.Plot(c, plt)
	}
	if pts.ybars != nil {
		pts.ybars.LineStyle.Color = pts.GlyphStyle.Color
		pts.ybars.Plot(c, plt)
	}
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger
// interface.
func (pts *S2D) DataRange() (xmin, xmax, ymin, ymax float64) {
	if dr, ok := pts.Data.(plot.DataRanger); ok {
		xmin, xmax, ymin, ymax = dr.DataRange()
	} else {
		xmin, xmax, ymin, ymax = plotter.XYRange(pts.Data)
	}

	if pts.xbars != nil {
		xmin1, xmax1, ymin1, ymax1 := pts.xbars.DataRange()
		xmin = math.Min(xmin1, xmin)
		xmax = math.Max(xmax1, xmax)
		ymin = math.Min(ymin1, ymin)
		ymax = math.Max(ymax1, ymax)
	}

	if pts.ybars != nil {
		xmin1, xmax1, ymin1, ymax1 := pts.ybars.DataRange()
		xmin = math.Min(xmin1, xmin)
		xmax = math.Max(xmax1, xmax)
		ymin = math.Min(ymin1, ymin)
		ymax = math.Max(ymax1, ymax)
	}

	return xmin, xmax, ymin, ymax
}

// GlyphBoxes returns a slice of plot.GlyphBoxes,
// implementing the plot.GlyphBoxer interface.
func (pts *S2D) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, pts.Data.Len())
	for i := 0; i < pts.Data.Len(); i++ {
		x, y := pts.Data.XY(i)
		bs[i].X = plt.X.Norm(x)
		bs[i].Y = plt.Y.Norm(y)
		bs[i].Rectangle = pts.GlyphStyle.Rectangle()
	}
	if pts.xbars != nil {
		bs = append(bs, pts.xbars.GlyphBoxes(plt)...)
	}
	if pts.ybars != nil {
		bs = append(bs, pts.ybars.GlyphBoxes(plt)...)
	}
	return bs
}

// Thumbnail the thumbnail for the Scatter,
// implementing the plot.Thumbnailer interface.
func (pts *S2D) Thumbnail(c *draw.Canvas) {
	c.DrawGlyph(pts.GlyphStyle, c.Center())
}
