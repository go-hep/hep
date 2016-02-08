// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package hplot

import (
	"math"

	"golang.org/x/exp/shiny/screen"

	"github.com/go-hep/hplot/vgshiny"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
)

// Plot is the basic type representing a plot.
type Plot struct {
	plot.Plot
}

// New returns a new plot with some reasonable
// default settings.
func New() (*Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	// p.X.Padding = 0
	// p.Y.Padding = 0
	// p.Style = GnuplotStyle{}
	return &Plot{*p}, nil
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

// Show displays the plot to the screen, with the given dimensions
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
