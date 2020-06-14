// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// BinnedErrBand implements the plot.Plotter interface,
// drawing a colored band for the error on any binned
// quantity.
type BinnedErrBand struct {

	// Error for each bins
	YErrs *plotter.YErrors

	// Definition of the bins.
	// FIX-ME[rmadar]: maybe better to use hbook.Range?
	Bins [][2]float64

	// LineStyle is the style of the line contouring the band.
	// Use zero width to disable.
	draw.LineStyle

	// FillColor is the color to fill the area between
	// the top and bottom data points.
	// Use nil to disable the filling.
	FillColor color.Color
}

// NewBinnedErrBand creates a binned error band from a binning
// (slice of range) and y errors bars.
// FIX-ME[rmadar]: use a more friendly type to pass error bars?
func NewBinnedErrBand(bins [][2]float64, yerrs *plotter.YErrors) *BinnedErrBand {
	return &BinnedErrBand{
		YErrs: yerrs,
		Bins:  bins,
	}
}


func (b *BinnedErrBand) Plot(c draw.Canvas, plt *plot.Plot) {

	for i := range b.Bins {

		// Get four corner of the ith bin
		xmin, xmax := b.Bins[i][0], b.Bins[i][1]
		ymin, ymax := b.YErrs.YError(i)		

		xys := plotter.XYs{
			plotter.XY{X: xmin, Y: ymin},
			plotter.XY{X: xmin, Y: ymax},
			plotter.XY{X: xmax, Y: ymax},
			plotter.XY{X: xmax, Y: ymin},
		}

		// Polygon
		// FIX-ME(rmadar): it would be better to draw only top
		//                 and bottom horizontal lines (not vertical lines)
		poly := plotter.Polygon{
			XYs:       []plotter.XYs{xys},
			LineStyle: b.LineStyle,
			Color:     b.FillColor,
		}
		poly.Plot(c, plt)
	}
}


func (band *BinnedErrBand) DataRange() (xmin, xmax, ymin, ymax float64) {
	return 
}
