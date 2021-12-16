// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"io"
	"math"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Plot is the basic type representing a plot.
type Plot struct {
	*plot.Plot
	Style Style
}

// muNewPlot protects access to gonum/plot.DefaultFont
var muNewPlot sync.Mutex

// New returns a new plot with some reasonable
// default settings.
func New() *Plot {
	muNewPlot.Lock()
	defer muNewPlot.Unlock()

	style := DefaultStyle
	defer style.reset(plot.DefaultFont)
	plot.DefaultFont = style.Fonts.Default

	pp := &Plot{
		Plot:  plot.New(),
		Style: style,
	}
	pp.Style.Apply(pp)
	return pp
}

// Add adds a Plotters to the plot.
//
// If the plotters implements DataRanger then the
// minimum and maximum values of the X and Y
// axes are changed if necessary to fit the range of
// the data.
//
// When drawing the plot, Plotters are drawn in the
// order in which they were added to the plot.
func (p *Plot) Add(ps ...plot.Plotter) {
	for _, d := range ps {
		if x, ok := d.(plot.DataRanger); ok {
			xmin, xmax, ymin, ymax := x.DataRange()
			p.Plot.X.Min = math.Min(p.Plot.X.Min, xmin)
			p.Plot.X.Max = math.Max(p.Plot.X.Max, xmax)
			p.Plot.Y.Min = math.Min(p.Plot.Y.Min, ymin)
			p.Plot.Y.Max = math.Max(p.Plot.Y.Max, ymax)
		}
	}

	p.Plot.Add(ps...)
}

// Save saves the plot to an image file.  The file format is determined
// by the extension.
//
// Supported extensions are:
//
//  .eps, .jpg, .jpeg, .pdf, .png, .svg, .tex, .tif and .tiff.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (p *Plot) Save(w, h vg.Length, file string) error {
	return Save(p, w, h, file)
}

// WriterTo returns an io.WriterTo that will write the plot as
// the specified image format.
//
// Supported formats are:
//
//  eps, jpg|jpeg, pdf, png, svg, tex and tif|tiff.
func (p *Plot) WriterTo(w, h vg.Length, format string) (io.WriterTo, error) {
	return WriterTo(p, w, h, format)
}

// Draw draws a plot to a draw.Canvas.
//
// Plotters are drawn in the order in which they were
// added to the plot.  Plotters that  implement the
// GlyphBoxer interface will have their GlyphBoxes
// taken into account when padding the plot so that
// none of their glyphs are clipped.
func (p *Plot) Draw(dc draw.Canvas) {
	p.Plot.Draw(dc)
}

var (
	_ Drawer = (*Plot)(nil)
)
