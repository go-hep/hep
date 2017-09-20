// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package hplot // import "go-hep.org/x/hep/hplot"

import (
	"math"

	"golang.org/x/exp/shiny/screen"

	"go-hep.org/x/hep/hplot/vgshiny"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// Plot is the basic type representing a plot.
type Plot struct {
	*plot.Plot
	Style Style
}

// New returns a new plot with some reasonable
// default settings.
func New() (*Plot, error) {
	style := DefaultStyle
	defer style.reset(plot.DefaultFont)
	plot.DefaultFont = style.Fonts.Name

	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	pp := &Plot{
		Plot:  p,
		Style: style,
	}
	pp.Style.Apply(pp)
	// p.X.Padding = 0
	// p.Y.Padding = 0
	// p.Style = GnuplotStyle{}
	return pp, nil
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
//  .eps, .jpg, .jpeg, .pdf, .png, .svg, .tif and .tiff.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (p *Plot) Save(w, h vg.Length, file string) (err error) {
	switch {
	case w <= 0 && h <= 0:
		w = vgimg.DefaultWidth
		h = vgimg.DefaultWidth / math.Phi
	case w <= 0:
		w = h * math.Phi
	case h <= 0:
		h = w / math.Phi
	}
	return p.Plot.Save(w, h, file)
}

// Show displays the plot to the screen, with the given dimensions.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func (p *Plot) Show(w, h vg.Length, scr screen.Screen) (*vgshiny.Canvas, error) {
	switch {
	case w <= 0 && h <= 0:
		w = vgimg.DefaultWidth
		h = vgimg.DefaultWidth / math.Phi
	case w <= 0:
		w = h * math.Phi
	case h <= 0:
		h = w / math.Phi
	}
	c, err := vgshiny.New(scr, w, h)
	if err != nil {
		return nil, err
	}
	p.Draw(draw.New(c))
	c.Paint()
	return c, err
}

// zip zips together 2 slices and implements the plotter.XYer interface.
type zip struct {
	x []float64
	y []float64
}

// Len implements the plotter.XYer interface
func (z zip) Len() int { return len(z.x) }

// XY implements the plotter.XYer interface
func (z zip) XY(i int) (x, y float64) { return z.x[i], z.y[i] }

// ZipXY zips together 2 slices x and y in such a way to implement the
// plotter.XYer interface.
//
// ZipXY panics if the slices are not of the same length.
func ZipXY(x, y []float64) plotter.XYer {
	if len(x) != len(y) {
		panic("hplot: slices length differ")
	}
	return zip{x: x, y: y}
}

// Options encodes various options to pass to a plot.
type Options int32

const (
	OptNone      Options = 1 << iota // default
	WithXErrBars                     // enable display of X-error bars (if available)
	WithYErrBars                     // enable display of Y-error bars (if available)
)
