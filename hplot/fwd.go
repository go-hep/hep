// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// NewLine returns a Line that uses the default line style and
// does not draw glyphs.
func NewLine(xys plotter.XYer) (*plotter.Line, error) {
	return plotter.NewLine(xys)
}

// NewScatter returns a Scatter that uses the
// default glyph style.
func NewScatter(xys plotter.XYer) (*plotter.Scatter, error) {
	return plotter.NewScatter(xys)
}

// NewGrid returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewGrid() *plotter.Grid {
	return plotter.NewGrid()
}

// A Legend gives a description of the meaning of different data elements of
// the plot. Each legend entry has a name and a thumbnail, where the thumbnail
// shows a small sample of the display style of the corresponding data.
type Legend = plot.Legend

// NewLegend returns a legend with the default parameter settings.
func NewLegend() Legend {
	return plot.NewLegend()
}
