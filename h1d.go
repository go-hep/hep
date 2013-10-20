package hplot

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/go-hep/hist"
	"github.com/go-hep/hplot/plotinum/plot"
	"github.com/go-hep/hplot/plotinum/plotter"
	"github.com/go-hep/hplot/plotinum/vg"
)

type h1d struct {
	hist.Hist1D
}

func (h h1d) XY(i int) (float64, float64) {
	axis := h.Axis()
	xmin := axis.BinLowerEdge(i)
	xmax := axis.BinUpperEdge(i)
	x := 0.5 * (xmax - xmin)
	y := h.Content(i)
	return x, y
}

func (h h1d) Len() int {
	return int(h.Axis().Bins())
}

func (h h1d) Value(idx int) float64 {
	return h.Content(idx)
}

// Histogram implements the plotter.Plotter interface,
// drawing a histogram of the data.
type Histogram struct {
	// Hist is the histogramming data
	Hist h1d

	// FillColor is the color used to fill each
	// bar of the histogram.  If the color is nil
	// then the bars are not filled.
	FillColor color.Color

	// LineStyle is the style of the outline of each
	// bar of the histogram.
	plot.LineStyle

	// InfoStyle is the style of infos displayed for
	// the histogram (entries, mean, rms)
	Infos HInfos
}

type HInfoStyle int

const (
	HInfo_DefaultStyle HInfoStyle = 0 // HInfo_Entries | HInfo_Mean | HInfo_RMS
	HInfo_Entries      HInfoStyle = iota << 1
	HInfo_Mean
	HInfo_RMS
)

type HInfos struct {
	Style HInfoStyle
}

// NewHistogram returns a new histogram
// that represents the distribution of values
// using the given number of bins.
//
// Each y value is assumed to be the frequency
// count for the corresponding x.
//
// If the number of bins is non-positive than
// a reasonable default is used.
func NewHistogram(xy plotter.XYer, n int) (*Histogram, error) {
	if n <= 0 {
		return nil, errors.New("Histogram with non-positive number of bins")
	}
	h := hist_from_xyer(xy, n)
	return NewH1D(h)
}

// NewHist returns a new histogram, as in
// NewHistogram, except that it accepts a plotter.Valuer
// instead of an XYer.
func NewHist(vs plotter.Valuer, n int) (*Histogram, error) {
	return NewHistogram(unitYs{vs}, n)
}

type unitYs struct {
	plotter.Valuer
}

func (u unitYs) XY(i int) (float64, float64) {
	return u.Value(i), 1.0
}

// NewH1D returns a new histogram, as in
// NewHistogram, except that it accepts a hist.Hist1D
// instead of a plotter.XYer
func NewH1D(h hist.Hist1D) (*Histogram, error) {
	return &Histogram{
		Hist:      h1d{h},
		FillColor: color.White,
		LineStyle: plotter.DefaultLineStyle,
	}, nil
}

// DataRange returns the minimum and maximum X and Y values
func (h *Histogram) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = h.Hist.Axis().LowerEdge()
	xmax = h.Hist.Axis().UpperEdge()
	ymin = h.Hist.Min()
	ymax = h.Hist.Max()
	return
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (h *Histogram) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	var pts []plot.Point
	hist := h.Hist
	axis := hist.Axis()
	nbins := int(axis.Bins())
	for bin := 0; bin < nbins; bin++ {
		switch bin {
		case 0:
			pts = append(pts, plot.Pt(trX(axis.BinLowerEdge(bin)), trY(0)))
			pts = append(pts, plot.Pt(trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))))
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))))

		case nbins - 1:
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin-1)), trY(hist.Content(bin-1))))
			pts = append(pts, plot.Pt(trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))))
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))))
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin)), trY(0.)))

		default:
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin-1)), trY(hist.Content(bin-1))))
			pts = append(pts, plot.Pt(trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))))
			pts = append(pts, plot.Pt(trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))))
		}
	}

	if h.FillColor != nil {
		da.FillPolygon(h.FillColor, da.ClipPolygonXY(pts))
	}
	da.StrokeLines(h.LineStyle, da.ClipLinesXY(pts)...)

	fnt, err := vg.MakeFont(plotter.DefaultFont, plotter.DefaultFontSize)
	if err == nil {
		sty := plot.TextStyle{Font: fnt}
		legend := hist_legend{
			ColWidth:  plotter.DefaultFontSize,
			TextStyle: sty,
		}

		switch h.Infos.Style {
		case HInfo_DefaultStyle:
			legend.Add("Entries", hist.Entries())
			legend.Add("Mean", hist.Mean())
			legend.Add("RMS", hist.RMS())
		case HInfo_Entries:
			legend.Add("Entries", hist.Entries())
		case HInfo_Mean:
			legend.Add("Mean", hist.Mean())
		case HInfo_RMS:
			legend.Add("RMS", hist.RMS())
		default:
		}
		legend.Top = true

		legend.draw(da)
	}
}

// GlyphBoxes returns a slice of GlyphBoxes,
// one for each of the bins, implementing the
// plot.GlyphBoxer interface.
func (h *Histogram) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	axis := h.Hist.Axis()
	bs := make([]plot.GlyphBox, axis.Bins())
	for i, _ := range bs {
		y := h.Hist.Content(i)
		xmin := axis.BinLowerEdge(i)
		w := p.X.Norm(axis.BinWidth(i))
		bs[i].X = p.X.Norm(xmin + 0.5*w)
		bs[i].Y = p.Y.Norm(y)
		//h := vg.Points(1e-5) //1 //p.Y.Norm(axis.BinWidth(i))
		bs[i].Rect.Min.X = vg.Length(xmin - 0.5*w)
		bs[i].Rect.Min.Y = vg.Length(y - 0.5*w)
		bs[i].Rect.Size.X = vg.Length(w)
		bs[i].Rect.Size.Y = vg.Length(0)

		r := vg.Points(5)
		//r = vg.Length(w)
		bs[i].Rect.Min = plot.Pt(0, 0)
		bs[i].Rect.Size = plot.Pt(0, r)
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

func hist_from_xyer(xys plotter.XYer, n int) h1d {
	xmin, xmax := plotter.Range(plotter.XValues{xys})
	h := hist.NewHist1D(n, xmin, xmax)

	for i := 0; i < xys.Len(); i++ {
		x, y := xys.XY(i)
		h.Fill(x, y)
	}

	return h1d{h}
}

// A Legend gives a description of the meaning of different
// data elements of the plot.  Each legend entry has a name
// and a thumbnail, where the thumbnail shows a small
// sample of the display style of the corresponding data.
type hist_legend struct {
	// TextStyle is the style given to the legend
	// entry texts.
	plot.TextStyle

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

// draw draws the legend to the given DrawArea.
func (l *hist_legend) draw(da plot.DrawArea) {
	textx := da.Min.X
	hdr := l.entryWidth() //+ l.TextStyle.Width(" ")
	l.ColWidth = hdr
	valx := textx + l.ColWidth + l.TextStyle.Width(" ")
	if !l.Left {
		textx = da.Max().X - l.ColWidth
		valx = textx - l.TextStyle.Width(" ")
	}
	valx += l.XOffs
	textx += l.XOffs

	enth := l.entryHeight()
	y := da.Max().Y - enth
	if !l.Top {
		y = da.Min.Y + (enth+l.Padding)*(vg.Length(len(l.entries))-1)
	}
	y += l.YOffs

	colx := &plot.DrawArea{
		Canvas: da.Canvas,
		Rect: plot.Rect{
			Min:  plot.Point{da.Min.X, y},
			Size: plot.Point{2 * l.ColWidth, enth},
		},
	}
	for _, e := range l.entries {
		yoffs := (enth - l.TextStyle.Height(e.text)) / 2
		da.FillText(l.TextStyle, textx-hdr, colx.Min.Y+yoffs, +0, 0, e.text)
		da.FillText(l.TextStyle, textx+hdr, colx.Min.Y+yoffs, -1, 0, e.value)
		colx.Min.Y -= enth + l.Padding
	}

	bbox_xmin := textx - hdr - l.TextStyle.Width(" ")
	bbox_xmax := da.Max().X
	bbox_ymin := colx.Min.Y + enth
	bbox_ymax := da.Max().Y
	bbox := []plot.Point{
		{bbox_xmin, bbox_ymax},
		{bbox_xmin, bbox_ymin},
		{bbox_xmax, bbox_ymin},
		{bbox_xmax, bbox_ymax},
		{bbox_xmin, bbox_ymax},
	}
	da.StrokeLines(plotter.DefaultLineStyle, bbox)
}

// entryHeight returns the height of the tallest legend
// entry text.
func (l *hist_legend) entryHeight() (height vg.Length) {
	for _, e := range l.entries {
		if h := l.TextStyle.Height(e.text); h > height {
			height = h
		}
	}
	return
}

// entryWidth returns the width of the largest legend
// entry text.
func (l *hist_legend) entryWidth() (width vg.Length) {
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
func (l *hist_legend) Add(name string, value interface{}) {
	str := ""
	switch value.(type) {
	case float64, float32:
		str = fmt.Sprintf("%6.4g ", value)
	default:
		str = fmt.Sprintf("%v ", value)
	}
	l.entries = append(l.entries, legendEntry{text: name, value: str})
}

// crop returns a new DrawArea corresponding to the receiver
// area with the given number of inches added to the minimum
// and maximum x and y values of the DrawArea's Rect.
func crop_da(da plot.DrawArea, minx, miny, maxx, maxy vg.Length) plot.DrawArea {
	minpt := plot.Point{
		X: da.Min.X + minx,
		Y: da.Min.Y + miny,
	}
	sz := plot.Point{
		X: da.Max().X + maxx - minpt.X,
		Y: da.Max().Y + maxy - minpt.Y,
	}
	return plot.DrawArea{
		Canvas: vg.Canvas(da),
		Rect:   plot.Rect{Min: minpt, Size: sz},
	}
}

// EOF
