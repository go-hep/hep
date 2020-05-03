// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// H1D implements the plotter.Plotter interface,
// drawing a histogram of the data.
type H1D struct {
	// Hist is the histogramming data
	Hist *hbook.H1D

	// FillColor is the color used to fill each
	// bar of the histogram.  If the color is nil
	// then the bars are not filled.
	FillColor color.Color

	// LineStyle is the style of the outline of each
	// bar of the histogram.
	draw.LineStyle

	// GlyphStyle is the style of the glyphs drawn
	// at the top of each histogram bar.
	GlyphStyle draw.GlyphStyle

	// LogY allows rendering with a log-scaled Y axis.
	// When enabled, histogram bins with no entries will be discarded from
	// the histogram's DataRange.
	// The lowest Y value for the DataRange will be corrected to leave an
	// arbitrary amount of height for the smallest bin entry so it is visible
	// on the final plot.
	LogY bool

	// InfoStyle is the style of infos displayed for
	// the histogram (entries, mean, rms)
	Infos HInfos

	// YErrs is the y error bars plotter.
	YErrs *plotter.YErrorBars

	// Band displays a colored band between the y-min and y-max error bars.
	Band *Band
}

type HInfoStyle uint32

const (
	HInfoNone    HInfoStyle = 0
	HInfoEntries HInfoStyle = 1 << iota
	HInfoMean
	HInfoRMS
	HInfoStdDev
	HInfoSummary HInfoStyle = HInfoEntries | HInfoMean | HInfoStdDev
)

type HInfos struct {
	Style HInfoStyle
}

// NewH1FromXYer returns a new histogram
// that represents the distribution of values
// using the given number of bins.
//
// Each y value is assumed to be the frequency
// count for the corresponding x.
//
// It panics if the number of bins is non-positive.
func NewH1FromXYer(xy plotter.XYer, n int, opts ...Options) *H1D {
	if n <= 0 {
		panic(errors.New("hplot: histogram with non-positive number of bins"))
	}
	h := newHistFromXYer(xy, n)
	return NewH1D(h, opts...)
}

// NewH1FromValuer returns a new histogram, as in
// NewH1FromXYer, except that it accepts a plotter.Valuer
// instead of an XYer.
func NewH1FromValuer(vs plotter.Valuer, n int, opts ...Options) *H1D {
	return NewH1FromXYer(unitYs{vs}, n, opts...)
}

type unitYs struct {
	plotter.Valuer
}

func (u unitYs) XY(i int) (float64, float64) {
	return u.Value(i), 1.0
}

// NewH1D returns a new histogram, as in
// NewH1DFromXYer, except that it accepts a hbook.H1D
// instead of a plotter.XYer
func NewH1D(h *hbook.H1D, opts ...Options) *H1D {
	h1 := &H1D{
		Hist:      h,
		LineStyle: plotter.DefaultLineStyle,
	}

	cfg := newConfig(opts)

	h1.LogY = cfg.log.y
	h1.Infos = cfg.hinfos

	if cfg.band {
		_ = h1.withBand()
	}

	if cfg.bars.yerrs {
		h1.YErrs = h1.withYErrBars(nil)
	}

	if cfg.glyph != (draw.GlyphStyle{}) {
		h1.GlyphStyle = cfg.glyph
	}

	return h1
}

// withYErrBars enables the Y error bars
func (h *H1D) withYErrBars(yoffs []float64) *plotter.YErrorBars {
	bins := h.Hist.Binning.Bins
	if yoffs == nil {
		yoffs = make([]float64, len(bins))
	}
	data := make(plotter.XYs, 0, len(bins))
	yerr := make(plotter.YErrors, 0, len(bins))
	for i, bin := range bins {
		if bin.Entries() == 0 {
			continue
		}
		data = append(data, plotter.XY{
			X: bin.XMid(),
			Y: yoffs[i] + bin.SumW(),
		})
		ey := 0.5 * bin.ErrW()
		yerr = append(yerr, struct{ Low, High float64 }{ey, ey})
	}

	type yerrT struct {
		plotter.XYer
		plotter.YErrorer
	}

	yplt, err := plotter.NewYErrorBars(yerrT{data, yerr})
	if err != nil {
		panic(err)
	}
	yplt.LineStyle.Color = h.LineStyle.Color
	yplt.LineStyle.Width = h.LineStyle.Width

	return yplt
}

// withBand enables the band between ymin-ymax error bars.
func (h1 *H1D) withBand() error {

	bins := h1.Hist.Binning.Bins
	var (
		top = make(plotter.XYs, 2*len(bins))
		bot = make(plotter.XYs, 2*len(bins))
	)

	for i := range top {
		ibin := i / 2
		bin := bins[ibin]
		xmin, xmax := bin.XEdges().Min, bin.XEdges().Max
		switch {
		case i%2 != 0:
			top[i].X = xmax
			top[i].Y = bin.SumW() - 0.5*bin.ErrW()
			bot[i].X = xmax
			bot[i].Y = bin.SumW() + 0.5*bin.ErrW()
		default:
			top[i].X = xmin
			top[i].Y = bin.SumW() - 0.5*bin.ErrW()
			bot[i].X = xmin
			bot[i].Y = bin.SumW() + 0.5*bin.ErrW()
		}
	}

	h1.Band = NewBand(color.Gray{200}, top, bot)
	return nil
}

// DataRange returns the minimum and maximum X and Y values
func (h *H1D) DataRange() (xmin, xmax, ymin, ymax float64) {
	if !h.LogY {
		xmin, xmax, ymin, ymax = h.Hist.DataRange()
		if h.YErrs != nil {
			xmin1, xmax1, ymin1, ymax1 := h.YErrs.DataRange()
			xmin = math.Min(xmin, xmin1)
			ymin = math.Min(ymin, ymin1)
			xmax = math.Max(xmax, xmax1)
			ymax = math.Max(ymax, ymax1)
		}
		return xmin, xmax, ymin, ymax
	}

	xmin = math.Inf(+1)
	xmax = math.Inf(-1)
	ymin = math.Inf(+1)
	ymax = math.Inf(-1)
	ylow := math.Inf(+1) // ylow will hold the smallest non-zero y value.
	for _, bin := range h.Hist.Binning.Bins {
		xmax = math.Max(bin.XMax(), xmax)
		xmin = math.Min(bin.XMin(), xmin)
		ymax = math.Max(bin.SumW(), ymax)
		ymin = math.Min(bin.SumW(), ymin)
		if bin.SumW() != 0 {
			ylow = math.Min(bin.SumW(), ylow)
		}
	}

	if ymin == 0 && !math.IsInf(ylow, +1) {
		// Reserve a bit of space for the smallest bin to be displayed still.
		ymin = ylow * 0.5
	}

	if h.YErrs != nil {
		xmin1, xmax1, ymin1, ymax1 := h.YErrs.DataRange()
		xmin = math.Min(xmin, xmin1)
		ymin = math.Min(ymin, ymin1)
		xmax = math.Max(xmax, xmax1)
		ymax = math.Min(ymax, ymax1)
	}

	return
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (h *H1D) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)
	var pts []vg.Point
	hist := h.Hist
	bins := h.Hist.Binning.Bins
	nbins := len(bins)

	yfct := func(sumw float64) (ymin, ymax vg.Length) {
		return trY(0), trY(sumw)
	}
	if h.LogY {
		yfct = func(sumw float64) (ymin, ymax vg.Length) {
			ymin = c.Min.Y
			ymax = c.Min.Y
			if 0 != sumw {
				ymax = trY(sumw)
			}
			return ymin, ymax
		}
	}

	var glyphs []vg.Point

	for i, bin := range bins {
		xmin := trX(bin.XMin())
		xmax := trX(bin.XMax())
		sumw := bin.SumW()
		ymin, ymax := yfct(sumw)
		switch i {
		case 0:
			pts = append(pts, vg.Point{X: xmin, Y: ymin})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})

		case nbins - 1:
			lft := bins[i-1]
			xlft := trX(lft.XMax())
			_, ylft := yfct(lft.SumW())
			pts = append(pts, vg.Point{X: xlft, Y: ylft})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymin})

		default:
			lft := bins[i-1]
			xlft := trX(lft.XMax())
			_, ylft := yfct(lft.SumW())
			pts = append(pts, vg.Point{X: xlft, Y: ylft})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})
		}

		if h.GlyphStyle.Radius != 0 {
			x := trX(bin.XMid())
			_, y := yfct(bin.SumW())
			// capture glyph location, to be drawn after
			// the histogram line, if any.
			glyphs = append(glyphs, vg.Point{X: x, Y: y})
		}
	}

	if h.FillColor != nil {
		c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
	}

	if h.Band != nil {
		h.Band.Plot(c, p)
	}

	c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)

	if h.YErrs != nil {
		h.YErrs.Plot(c, p)
	}

	if h.GlyphStyle.Radius != 0 {
		for _, glyph := range glyphs {
			c.DrawGlyph(h.GlyphStyle, glyph)
		}
	}

	if h.Infos.Style != HInfoNone {
		fnt, err := vg.MakeFont(DefaultStyle.Fonts.Name, DefaultStyle.Fonts.Tick.Size)
		if err == nil {
			sty := draw.TextStyle{Font: fnt}
			legend := histLegend{
				ColWidth:  DefaultStyle.Fonts.Tick.Size,
				TextStyle: sty,
			}

			for i := uint32(0); i < 32; i++ {
				switch h.Infos.Style & (1 << i) {
				case HInfoEntries:
					legend.Add("Entries", hist.Entries())
				case HInfoMean:
					legend.Add("Mean", hist.XMean())
				case HInfoRMS:
					legend.Add("RMS", hist.XRMS())
				case HInfoStdDev:
					legend.Add("Std Dev", hist.XStdDev())
				default:
				}
			}
			legend.Top = true

			legend.draw(c)
		}
	}
}

// GlyphBoxes returns a slice of GlyphBoxes,
// one for each of the bins, implementing the
// plot.GlyphBoxer interface.
func (h *H1D) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	bins := h.Hist.Binning.Bins
	bs := make([]plot.GlyphBox, 0, len(bins))
	for i := range bins {
		bin := bins[i]
		y := bin.SumW()
		if h.LogY && y == 0 {
			continue
		}
		var box plot.GlyphBox
		xmin := bin.XMin()
		w := p.X.Norm(bin.XWidth())
		box.X = p.X.Norm(xmin + 0.5*w)
		box.Y = p.Y.Norm(y)
		box.Rectangle.Min.X = vg.Length(xmin - 0.5*w)
		box.Rectangle.Min.Y = vg.Length(y - 0.5*w)
		box.Rectangle.Max.X = vg.Length(w)
		box.Rectangle.Max.Y = vg.Length(0)

		r := vg.Points(5)
		box.Rectangle.Min = vg.Point{X: 0, Y: 0}
		box.Rectangle.Max = vg.Point{X: 0, Y: r}
		bs = append(bs, box)
	}
	return bs
}

// Normalize normalizes the histogram so that the
// total area beneath it sums to a given value.
// func (h *Histogram) Normalize(sum float64) {
// 	mass := 0.0
// 	for _, b := range h.Bins {
// 		mass += b.Weight
// 	}
// 	for i := range h.Bins {
// 		h.Bins[i].Weight *= sum / (h.Width * mass)
// 	}
// }

// Thumbnail draws a rectangle in the given style of the histogram.
func (h *H1D) Thumbnail(c *draw.Canvas) {
	ymin := c.Min.Y
	ymax := c.Max.Y
	xmin := c.Min.X
	xmax := c.Max.X
	dy := ymax - ymin

	// Style of the histogram
	hasFill := h.FillColor != nil
	hasLine := h.LineStyle.Width != 0
	hasGlyph := h.GlyphStyle != (draw.GlyphStyle{})
	hasBand := h.Band != nil

	// WIP [rmadar]: define default behaviour
	drawFill := hasFill
	drawBand := hasBand
	drawLine := hasLine
	drawGlyph := hasGlyph
	
	if drawFill {
		pts := []vg.Point{
			{X: xmin, Y: ymin},
			{X: xmax, Y: ymin},
			{X: xmax, Y: ymax},
			{X: xmin, Y: ymax},
			{X: xmin, Y: ymin},
		}
		c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
	}

	if drawBand {
		pts := []vg.Point{
			{X: xmin, Y: ymin + 0.2*dy},
			{X: xmax, Y: ymin + 0.2*dy},
			{X: xmax, Y: ymax - 0.2*dy},
			{X: xmin, Y: ymax - 0.2*dy},
			{X: xmin, Y: ymin + 0.2*dy},
		}
		c.FillPolygon(h.Band.FillColor, c.ClipPolygonXY(pts))
	}

	if drawLine {
		if hasFill && !hasGlyph && !hasBand {
			line := []vg.Point{
				{X: xmin, Y: ymin},
				{X: xmax, Y: ymin},
				{X: xmax, Y: ymax},
				{X: xmin, Y: ymax},
				{X: xmin, Y: ymin},
			}
			c.StrokeLines(h.LineStyle, c.ClipLinesX(line)...)
		} else {
			ymid := c.Center().Y
			line := []vg.Point{{X: xmin, Y: ymid}, {X: xmax, Y: ymid}}
			c.StrokeLines(h.LineStyle, c.ClipLinesX(line)...)
		}

	}

	if drawGlyph {
		c.DrawGlyph(h.GlyphStyle, c.Center())
		if h.YErrs != nil {
			var (
				yerrs = h.YErrs
				vsize = 0.5 * dy * 0.95
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

func newHistFromXYer(xys plotter.XYer, n int) *hbook.H1D {
	xmin, xmax := plotter.Range(plotter.XValues{XYer: xys})
	h := hbook.NewH1D(n, xmin, xmax)

	for i := 0; i < xys.Len(); i++ {
		x, y := xys.XY(i)
		h.Fill(x, y)
	}

	return h
}

// A Legend gives a description of the meaning of different
// data elements of the plot.  Each legend entry has a name
// and a thumbnail, where the thumbnail shows a small
// sample of the display style of the corresponding data.
type histLegend struct {
	// TextStyle is the style given to the legend
	// entry texts.
	draw.TextStyle

	// Padding is the amount of padding to add
	// betweeneach entry of the legend.  If Padding
	// is zero then entries are spaced based on the
	// font size.
	Padding vg.Length

	// Top and Left specify the location of the legend.
	// If Top is true the legend is located along the top
	// edge of the plot, otherwise it is located along
	// the bottom edge.  If Left is true then the legend
	// is located along the left edge of the plot, and the
	// text is positioned after the icons, otherwise it is
	// located along the right edge and the text is
	// positioned before the icons.
	Top, Left bool

	// XOffs and YOffs are added to the legend's
	// final position.
	XOffs, YOffs vg.Length

	// ColWidth is the width of legend names
	ColWidth vg.Length

	// entries are all of the legendEntries described
	// by this legend.
	entries []legendEntry
}

// A legendEntry represents a single line of a legend, it
// has a name and an icon.
type legendEntry struct {
	// text is the text associated with this entry.
	text string

	// value is the value associated with this entry
	value string
}

// draw draws the legend to the given canvas.
func (l *histLegend) draw(c draw.Canvas) {
	textx := c.Min.X
	hdr := l.entryWidth() //+ l.TextStyle.Width(" ")
	l.ColWidth = hdr
	if !l.Left {
		textx = c.Max.X - l.ColWidth
	}
	textx += l.XOffs

	enth := l.entryHeight()
	y := c.Max.Y - enth
	if !l.Top {
		y = c.Min.Y + (enth+l.Padding)*(vg.Length(len(l.entries))-1)
	}
	y += l.YOffs

	colx := &draw.Canvas{
		Canvas: c.Canvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: c.Min.X, Y: y},
			Max: vg.Point{X: 2 * l.ColWidth, Y: enth},
		},
	}
	for _, e := range l.entries {
		yoffs := (enth - l.TextStyle.Height(e.text)) / 2
		txt := l.TextStyle
		txt.XAlign = draw.XLeft
		c.FillText(txt, vg.Point{X: textx - hdr, Y: colx.Min.Y + yoffs}, e.text)
		txt.XAlign = draw.XRight
		c.FillText(txt, vg.Point{X: textx + hdr, Y: colx.Min.Y + yoffs}, e.value)
		colx.Min.Y -= enth + l.Padding
	}

	bboxXmin := textx - hdr - l.TextStyle.Width(" ")
	bboxXmax := c.Max.X
	bboxYmin := colx.Min.Y + enth
	bboxYmax := c.Max.Y
	bbox := []vg.Point{
		{X: bboxXmin, Y: bboxYmax},
		{X: bboxXmin, Y: bboxYmin},
		{X: bboxXmax, Y: bboxYmin},
		{X: bboxXmax, Y: bboxYmax},
		{X: bboxXmin, Y: bboxYmax},
	}
	c.StrokeLines(plotter.DefaultLineStyle, bbox)
}

// entryHeight returns the height of the tallest legend
// entry text.
func (l *histLegend) entryHeight() (height vg.Length) {
	for _, e := range l.entries {
		if h := l.TextStyle.Height(e.text); h > height {
			height = h
		}
	}
	return
}

// entryWidth returns the width of the largest legend
// entry text.
func (l *histLegend) entryWidth() (width vg.Length) {
	for _, e := range l.entries {
		if w := l.TextStyle.Width(e.value); w > width {
			width = w
		}
	}
	return
}

// Add adds an entry to the legend with the given name.
// The entry's thumbnail is drawn as the composite of all of the
// thumbnails.
func (l *histLegend) Add(name string, value interface{}) {
	str := ""
	switch value.(type) {
	case float64, float32:
		str = fmt.Sprintf("%6.4g ", value)
	default:
		str = fmt.Sprintf("%v ", value)
	}
	l.entries = append(l.entries, legendEntry{text: name, value: str})
}

var (
	_ plot.Plotter     = (*H1D)(nil)
	_ plot.Thumbnailer = (*H1D)(nil)
)
