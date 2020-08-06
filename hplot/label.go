// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Label struct {
	Text      string         // Text of the label
	X, Y      float64        // Position of the label
	TextStyle draw.TextStyle // Text style of the label

	// Normalized indicates whether the label position
	// is in data coordinates or normalized with regard
	// to the canvas space.
	// When normalized, the label position is assumed
	// to fall in the [0, 1] interval. If true, NewLabel
	// panics if x or y are outside [0, 1].
	// Normalized is false by default.
	Normalized bool

	// AutoAdjust enables auto adjustment of the label
	// position, when Normalized is true and when x
	// and/or y are close to 1 and the label is partly
	// outside the canvas. If false and the label doesn't
	// fit in the canvas, an error is returned.
	// AutoAdjust is false by default.
	AutoAdjust bool
}

// NewLabel creates a new label value from x, y.
func NewLabel(x, y float64, txt string, opts ...LabelOption) Label {

	// FIXME[rmadar]: default style to be handle properly
	//                once hplot.Style is ready.
	font, fSize := plotter.DefaultFont, plotter.DefaultFontSize
	defaultFont, err := vg.MakeFont(font, fSize)
	if err != nil {
		panic(fmt.Errorf("hplot: could not create font (%s, %v): %w", font, fSize, err))
	}
	style := draw.TextStyle{
		Color: color.Black,
		Font:  defaultFont,
	}

	// Default config
	cfg := &labelConfig{
		TextStyle:  style,
		Normalized: false,
		AutoAdjust: false,
	}

	// User-defined options
	for _, opt := range opts {
		opt(cfg)
	}

	// Sanity check
	if cfg.Normalized && (x < 0 || x > 1 || y < 0 || y > 1) {
		panic("hplot: normalized label position is outside [0,1]")
	}

	// Return the configured Label
	return Label{
		Text:       txt,
		X:          x,
		Y:          y,
		TextStyle:  cfg.TextStyle,
		Normalized: cfg.Normalized,
		AutoAdjust: cfg.Normalized,
	}

}

// Plot implements the Plotter interface,
// drawing the label on the canvas.
func (lbl Label) Plot(c draw.Canvas, p *plot.Plot) {
	lbls := lbl.labels(c, p)
	lbls.Plot(c, p)
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (lbl Label) DataRange() (xmin, xmax, ymin, ymax float64) {
	if lbl.Normalized {
		return math.Inf(+1), math.Inf(-1), math.Inf(+1), math.Inf(-1)
	}

	pLabels := lbl.labels(draw.Canvas{}, nil)
	return pLabels.DataRange()
}

// GlyphBoxes returns a GlyphBoxe, corresponding
// to the label, implementing the plot.GlyphBoxer interface.
func (lbl Label) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	if lbl.Normalized {
		return []plot.GlyphBox{
			{X: lbl.X, Y: lbl.Y, Rectangle: lbl.TextStyle.Rectangle(lbl.Text)},
		}
	}
	return lbl.labels(draw.Canvas{}, p).GlyphBoxes(p)
}

// Internal helper function to get plotter.Labels type.
func (lbl *Label) labels(c draw.Canvas, p *plot.Plot) *plotter.Labels {

	x := lbl.X
	y := lbl.Y
	if lbl.Normalized {

		// Check wether the label fits in the canvas
		txtBox := lbl.TextStyle.Rectangle(lbl.Text)
		xMax := p.X.Norm(float64(txtBox.Max.X))
		yMax := p.Y.Norm(float64(txtBox.Max.Y))
		if xMax > 1 || yMax > 1 {
			//if lbl.AutoAdjust {
			//	x = x - (xMax - 1)
			//	y = y - (yMax - 1)
			//} else {
			//panic("hplot: labels fall outside canvas")
			//}
		}

		// Turn relative into absolute coordinates
		dataCoord := func(xrel, xmin, xmax float64) float64 {
			return xmin + xrel*(xmax-xmin)
		}
		x = dataCoord(x, p.X.Min, p.X.Max)
		y = dataCoord(y, p.Y.Min, p.Y.Max)
	}

	xyL := plotter.XYLabels{
		XYs:    []plotter.XY{{X: x, Y: y}},
		Labels: []string{lbl.Text},
	}

	lbls, err := plotter.NewLabels(xyL)
	if err != nil {
		panic(fmt.Errorf("hplot: could not create labels: %w", err))
	}

	lbls.TextStyle = []draw.TextStyle{lbl.TextStyle}

	return lbls
}

type labelConfig struct {
	TextStyle  draw.TextStyle
	Normalized bool
	AutoAdjust bool
}

// Label option handles various options to
// configure the label.
type LabelOption func(cfg *labelConfig)

// WithTextStyle specifies the text style of the label.
func WithTextStyle(style draw.TextStyle) LabelOption {
	return func(cfg *labelConfig) {
		cfg.TextStyle = style
	}
}

// WithTextStyle specifies whether the coordinate are
// normalized to the canvas size.
func WithNormalized(norm bool) LabelOption {
	return func(cfg *labelConfig) {
		cfg.Normalized = norm
	}
}

// WithTextStyle specifies whether the coordinate are
// normalized to the canvas size.
func WithAutoAdjust(auto bool) LabelOption {
	return func(cfg *labelConfig) {
		cfg.AutoAdjust = auto
	}
}

var (
	_ plot.Plotter    = (*Label)(nil)
	_ plot.DataRanger = (*Label)(nil)
	_ plot.GlyphBoxer = (*Label)(nil)
)
