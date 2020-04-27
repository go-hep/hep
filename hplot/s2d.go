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

	XErrs *plotter.XErrorBars
	YErrs *plotter.YErrorBars

	// Band displays a colored band between the y-min and y-max error bars.
	Band *Band

	// Steps controls the style of the connecting
	// line (NoSteps, HiSteps, etc...)
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

	pts.XErrs = xplt
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

	pts.YErrs = yplt
	return nil
}

// withBand enables the band between ymin-ymax error bars.
func (pts *S2D) withBand() error {
	yerr, ok := pts.Data.(plotter.YErrorer)
	if !ok {
		return nil
	}

	var (
		top plotter.XYs
		bot plotter.XYs
	)

	switch pts.Steps {

	case NoSteps:
		top = make(plotter.XYs, pts.Data.Len())
		bot = make(plotter.XYs, pts.Data.Len())
		for i := range top {
			x, y := pts.Data.XY(i)
			ymin, ymax := yerr.YError(i)
			top[i].X = x
			top[i].Y = y + math.Abs(ymax)
			bot[i].X = x
			bot[i].Y = y - math.Abs(ymin)
		}

	case HiSteps:
		top = make(plotter.XYs, 2*pts.Data.Len())
		bot = make(plotter.XYs, 2*pts.Data.Len())
		xerr := pts.Data.(plotter.XErrorer)
		for i := range top {
			idata := i / 2
			x, y := pts.Data.XY(idata)
			xmin, xmax := xerr.XError(idata)
			ymin, ymax := yerr.YError(idata)
			switch {
			case i%2 != 0:
				top[i].X = x + math.Abs(xmax)
				top[i].Y = y + math.Abs(ymax)
				bot[i].X = x + math.Abs(xmax)
				bot[i].Y = y - math.Abs(ymin)
			default:
				top[i].X = x - math.Abs(xmin)
				top[i].Y = y + math.Abs(ymax)
				bot[i].X = x - math.Abs(xmin)
				bot[i].Y = y - math.Abs(ymin)
			}
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

		if pts.Steps == HiSteps {
			xerr := pts.Data.(plotter.XErrorer)
			dsteps := make(plotter.XYs, 0, 2*len(data))
			for i, d := range data {
				xmin, xmax := xerr.XError(i)
				dsteps = append(dsteps, plotter.XY{X: d.X - xmin, Y: d.Y})
				dsteps = append(dsteps, plotter.XY{X: d.X + xmax, Y: d.Y})
			}
			data = dsteps
		}

		line := plotter.Line{
			XYs:       data,
			LineStyle: pts.LineStyle,
		}
		line.Plot(c, plt)
	}

	if pts.XErrs != nil {
		pts.XErrs.LineStyle.Color = pts.GlyphStyle.Color
		pts.XErrs.Plot(c, plt)
	}
	if pts.YErrs != nil {
		pts.YErrs.LineStyle.Color = pts.GlyphStyle.Color
		pts.YErrs.Plot(c, plt)
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

	if pts.XErrs != nil {
		xmin1, xmax1, ymin1, ymax1 := pts.XErrs.DataRange()
		xmin = math.Min(xmin1, xmin)
		xmax = math.Max(xmax1, xmax)
		ymin = math.Min(ymin1, ymin)
		ymax = math.Max(ymax1, ymax)
	}

	if pts.YErrs != nil {
		xmin1, xmax1, ymin1, ymax1 := pts.YErrs.DataRange()
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
	if pts.XErrs != nil {
		bs = append(bs, pts.XErrs.GlyphBoxes(plt)...)
	}
	if pts.YErrs != nil {
		bs = append(bs, pts.YErrs.GlyphBoxes(plt)...)
	}
	return bs
}

// Thumbnail the thumbnail for the Scatter,
// implementing the plot.Thumbnailer interface.
func (pts *S2D) Thumbnail(c *draw.Canvas) {
	ymin := c.Min.Y
	ymax := c.Max.Y
	xmin := c.Min.X
	xmax := c.Max.X

	if pts.Band != nil {
		box := []vg.Point{
			{X: xmin, Y: ymin},
			{X: xmax, Y: ymin},
			{X: xmax, Y: ymax},
			{X: xmin, Y: ymax},
			{X: xmin, Y: ymin},
		}
		c.FillPolygon(pts.Band.FillColor, c.ClipPolygonXY(box))
	}

	if pts.LineStyle.Width != 0 {
		ymid := c.Center().Y
		line := []vg.Point{{X: xmin, Y: ymid}, {X: xmax, Y: ymid}}
		c.StrokeLines(pts.LineStyle, c.ClipLinesX(line)...)

	}

	if pts.GlyphStyle != (draw.GlyphStyle{}) {
		c.DrawGlyph(pts.GlyphStyle, c.Center())
		if pts.YErrs != nil {
			var (
				yerrs = pts.YErrs
				vsize = 0.5 * ((ymax - ymin) * 0.95)
				x     = c.Center().X
				ylo   = c.Center().Y - vsize
				yup   = c.Center().Y + vsize
				xylo  = vg.Point{X: x, Y: ylo}
				xyup  = vg.Point{X: x, Y: yup}
				line  = []vg.Point{xylo, xyup}
				bar   = c.ClipLinesY(line)
			)
			c.StrokeLines(yerrs.LineStyle, bar...)
			for _, pt := range []vg.Point{xylo, xyup} {
				if c.Contains(pt) {
					c.StrokeLine2(yerrs.LineStyle,
						pt.X-yerrs.CapWidth/2,
						pt.Y,
						pt.X+yerrs.CapWidth/2,
						pt.Y,
					)
				}
			}
		}
	}
}
