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
	"gonum.org/v1/plot/vg/draw"
)

// Label displays a user-defined text string on a plot.
//
// Fields of Label should not be modified once the Label has been
// added to an hplot.Plot.
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
	//
	// Normalized is false by default.
	Normalized bool

	// AutoAdjust enables auto adjustment of the label
	// position, when Normalized is true and when x
	// and/or y are close to 1 and the label is partly
	// outside the canvas. If false and the label doesn't
	// fit in the canvas, an error is returned.
	//
	// AutoAdjust is false by default.
	AutoAdjust bool

	// cache of gonum/plot.Labels
	plt *plotter.Labels
}

// NewLabel creates a new txt label at position (x, y).
func NewLabel(x, y float64, txt string, opts ...LabelOption) *Label {

	style := draw.TextStyle{
		Color:   color.Black,
		Font:    DefaultStyle.Fonts.Tick, // FIXME(sbinet): add a field in Style?
		Handler: DefaultStyle.TextHandler,
	}

	cfg := &labelConfig{
		TextStyle: style,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.TextStyle.Handler == nil {
		cfg.TextStyle.Handler = style.Handler
	}

	if cfg.Normalized {
		if !(0 <= x && x <= 1) {
			panic(fmt.Errorf(
				"hplot: normalized label x-position is outside [0,1]: %g", x,
			))
		}
		if !(0 <= y && y <= 1) {
			panic(fmt.Errorf(
				"hplot: normalized label y-position is outside [0,1]: %g", y,
			))
		}
	}

	return &Label{
		Text:       txt,
		X:          x,
		Y:          y,
		TextStyle:  cfg.TextStyle,
		Normalized: cfg.Normalized,
		AutoAdjust: cfg.AutoAdjust,
	}
}

// Plot implements the Plotter interface,
// drawing the label on the canvas.
func (lbl *Label) Plot(c draw.Canvas, p *plot.Plot) {
	lbl.labels(c, p).Plot(c, p)
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (lbl *Label) DataRange() (xmin, xmax, ymin, ymax float64) {
	if lbl.Normalized {
		return math.Inf(+1), math.Inf(-1), math.Inf(+1), math.Inf(-1)
	}

	return lbl.labels(draw.Canvas{}, nil).DataRange()
}

// GlyphBoxes returns a GlyphBox, corresponding to the label.
// GlyphBoxes implements the plot.GlyphBoxer interface.
func (lbl *Label) GlyphBoxes(p *plot.Plot) []plot.GlyphBox {
	if lbl.plt == nil {
		return nil
	}
	// we expect Label.Plot(c,p) has already been called.
	return lbl.labels(draw.Canvas{}, p).GlyphBoxes(p)
}

// Internal helper function to get plotter.Labels type.
func (lbl *Label) labels(c draw.Canvas, p *plot.Plot) *plotter.Labels {
	if lbl.plt != nil {
		return lbl.plt
	}

	var (
		x = lbl.X
		y = lbl.Y

		err error
	)

	if lbl.Normalized {
		// Check whether the label fits in the canvas
		box := lbl.TextStyle.Rectangle(lbl.Text)
		rect := c.Rectangle.Size()
		xmax := lbl.X + box.Max.X.Points()/rect.X.Points()
		ymax := lbl.Y + box.Max.Y.Points()/rect.Y.Points()
		if xmax > 1 || ymax > 1 {
			switch {
			case lbl.AutoAdjust:
				x, y = lbl.adjust(1/rect.X.Points(), 1/rect.Y.Points())
			default:
				panic(fmt.Errorf(
					"hplot: label (%g, %g) falls outside data canvas",
					x, y,
				))
			}
		}

		// Turn relative into absolute coordinates
		x = lbl.scale(x, p.X.Min, p.X.Max, p.X.Scale)
		y = lbl.scale(y, p.Y.Min, p.Y.Max, p.Y.Scale)
	}

	lbl.plt, err = plotter.NewLabels(plotter.XYLabels{
		XYs:    []plotter.XY{{X: x, Y: y}},
		Labels: []string{lbl.Text},
	})
	if err != nil {
		panic(fmt.Errorf("hplot: could not create labels: %w", err))
	}

	lbl.plt.TextStyle = []draw.TextStyle{lbl.TextStyle}

	return lbl.plt
}

func (lbl *Label) adjust(xnorm, ynorm float64) (x, y float64) {
	x = lbl.adjustX(xnorm)
	y = lbl.adjustY(ynorm)
	return x, y
}

func (lbl *Label) adjustX(xnorm float64) float64 {
	var (
		box  = lbl.TextStyle.Rectangle(lbl.Text)
		size = box.Size().X.Points() * xnorm
		x    = lbl.X
		dx   = size - (1 - x)
	)
	if x+size > 1 {
		x -= dx
	}
	if x < 0 {
		x = 0
	}
	return x
}

func (lbl *Label) adjustY(ynorm float64) float64 {
	var (
		box  = lbl.TextStyle.Rectangle(lbl.Text)
		size = box.Size().Y.Points() * ynorm
		y    = lbl.Y
		dy   = size - (1 - y)
	)
	if y+size > 1 {
		y -= dy
	}
	if y < 0 {
		y = 0
	}
	return y
}

func (Label) scale(v, min, max float64, scaler plot.Normalizer) float64 {
	mid := min + 0.5*(max-min)
	if math.Abs(scaler.Normalize(min, max, mid)-0.5) < 1e-12 {
		return min + v*(max-min)
	}

	// log-scale
	min = math.Log(min)
	max = math.Log(max)
	return math.Exp(min + v*(max-min))
}

type labelConfig struct {
	TextStyle  draw.TextStyle
	Normalized bool
	AutoAdjust bool
}

// LabelOption handles various options to configure a Label.
type LabelOption func(cfg *labelConfig)

// WithLabelTextStyle specifies the text style of the label.
func WithLabelTextStyle(style draw.TextStyle) LabelOption {
	return func(cfg *labelConfig) {
		cfg.TextStyle = style
	}
}

// WithLabelNormalized specifies whether the coordinates are
// normalized to the canvas size.
func WithLabelNormalized(norm bool) LabelOption {
	return func(cfg *labelConfig) {
		cfg.Normalized = norm
	}
}

// WithLabelAutoAdjust specifies whether the coordinates are
// automatically adjusted to the canvas size.
func WithLabelAutoAdjust(auto bool) LabelOption {
	return func(cfg *labelConfig) {
		cfg.AutoAdjust = auto
	}
}

var (
	_ plot.Plotter    = (*Label)(nil)
	_ plot.DataRanger = (*Label)(nil)
	_ plot.GlyphBoxer = (*Label)(nil)
)
