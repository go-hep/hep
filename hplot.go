// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package hplot

import (
	"math"

	"github.com/go-hep/hplot/plotinum/plot"
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
	p.X.Padding = 0
	p.Y.Padding = 0
	p.Style = plot.GnuplotStyle
	return &Plot{*p}, nil
}

func padf(xmin, xmax float64) (float64, float64) {
	if xmin < 0 && xmax < 0 {
		return xmin * 1.05, xmax * 0.95
	}
	if xmin < 0 && xmax >= 0 {
		return xmin * 1.05, xmax * 1.05
	}
	return xmin * 0.95, xmax * 1.05
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
			xmin, xmax = padf(xmin, xmax)
			ymin, ymax = padf(ymin, ymax)
			p.Plot.X.Min = math.Min(p.Plot.X.Min, xmin)
			p.Plot.X.Max = math.Max(p.Plot.X.Max, xmax)
			p.Plot.Y.Min = math.Min(p.Plot.Y.Min, ymin)
			p.Plot.Y.Max = math.Max(p.Plot.Y.Max, ymax)
		}
	}

	p.Plot.Add(ps...)
}

// EOF
