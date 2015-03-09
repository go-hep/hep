package main

import (
	"image/color"

	"github.com/go-hep/hbook"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

type h1d struct {
	*hbook.H1D
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
	draw.LineStyle
}

// NewH1D returns a new histogram, as in
// NewHistogram, except that it accepts a hist.Hist1D
// instead of a plotter.XYer
func NewH1D(h *hbook.H1D) (*Histogram, error) {
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
func (h *Histogram) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)
	var pts []draw.Point
	hist := h.Hist
	axis := hist.Axis()
	nbins := int(axis.Bins())
	for bin := 0; bin < nbins; bin++ {
		switch bin {
		case 0:
			pts = append(pts, draw.Point{trX(axis.BinLowerEdge(bin)), trY(0)})
			pts = append(pts, draw.Point{trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))})
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))})

		case nbins - 1:
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin - 1)), trY(hist.Content(bin - 1))})
			pts = append(pts, draw.Point{trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))})
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))})
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin)), trY(0.)})

		default:
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin - 1)), trY(hist.Content(bin - 1))})
			pts = append(pts, draw.Point{trX(axis.BinLowerEdge(bin)), trY(hist.Content(bin))})
			pts = append(pts, draw.Point{trX(axis.BinUpperEdge(bin)), trY(hist.Content(bin))})
		}
	}

	if h.FillColor != nil {
		c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
	}
	c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)
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
		bs[i].Rectangle.Min.X = vg.Length(xmin - 0.5*w)
		bs[i].Rectangle.Max.X = vg.Length(xmin + 0.5*w)
		bs[i].Rectangle.Min.Y = vg.Length(y - 0.5*w)
		bs[i].Rectangle.Max.Y = vg.Length(y + 0.5*w)

		r := vg.Points(5)
		//r = vg.Length(w)
		bs[i].Rectangle.Min = draw.Point{0, 0}
		bs[i].Rectangle.Max = draw.Point{0, r}
	}
	return bs
}
