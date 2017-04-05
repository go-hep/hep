// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// S2D plots a set of 2-dim points with error bars.
type S2D struct {
	Data plotter.XYer

	// GlyphStyle is the style of the glyphs drawn
	// at each point.
	draw.GlyphStyle

	// options controls various display options of this plot.
	options Options

	xbars *plotter.XErrorBars
	ybars *plotter.YErrorBars
}

// withXErrBars enables the X error bars
func (pts *S2D) withXErrBars() error {
	if pts.options&WithXErrBars == 0 {
		pts.xbars = nil
		return nil
	}
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
	if pts.options&WithYErrBars == 0 {
		pts.ybars = nil
		return nil
	}
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

// NewS2D creates a 2-dim scatter plot from a XYer.
func NewS2D(data plotter.XYer, opts ...Options) *S2D {
	s := &S2D{
		Data:       data,
		GlyphStyle: plotter.DefaultGlyphStyle,
	}
	for _, opt := range opts {
		s.options |= opt
	}

	for _, f := range []func() error{s.withXErrBars, s.withYErrBars} {
		_ = f()
	}

	s.GlyphStyle.Shape = draw.CrossGlyph{}
	return s
}

// Plot draws the Scatter, implementing the plot.Plotter
// interface.
func (pts *S2D) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	for i := 0; i < pts.Data.Len(); i++ {
		x, y := pts.Data.XY(i)
		c.DrawGlyph(pts.GlyphStyle, vg.Point{X: trX(x), Y: trY(y)})
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
