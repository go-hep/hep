// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"	
)

// BinnedErrBand implements the plot.Plotter interface,
// drawing a colored band for the error on any binned
// quantity.
type BinnedErrBand struct {

	// Y value for each bin
	Ys plotter.Values

	// Y error for each bins
	YErrs plotter.YErrors

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

// NewBinnedErrBand creates a binned error band
// from a binning (slice of range) and y errors bars.
// FIX-ME[rmadar]: use a more friendly type to pass Y errors?
func NewBinnedErrBand(bins [][2]float64, ys plotter.Values, yerrs plotter.YErrors) *BinnedErrBand {

	cpy, err := plotter.CopyValues(ys)
	if err != nil {
		log.Fatalf("cannot copy values")
	}

	return &BinnedErrBand{
		Ys:    cpy,
		YErrs: yerrs,
		Bins:  bins,
	}
}

// Plot implements the Plotter interface,
// drawing a colored box defined by width
// of bins (x-axis) and error (y-axis).
func (b *BinnedErrBand) Plot(c draw.Canvas, plt *plot.Plot) {

	for i, y := range b.Ys {

		// Get four corner of the ith bin
		xmin, xmax := b.Bins[i][0], b.Bins[i][1]
		ydo, yup := b.YErrs.YError(i)
		xys := plotter.XYs{
			plotter.XY{X: xmin, Y: y - ydo},
			plotter.XY{X: xmin, Y: y + yup},
			plotter.XY{X: xmax, Y: y + yup},
			plotter.XY{X: xmax, Y: y - ydo},
		}

		// Polygon
		poly := plotter.Polygon{
			XYs:   []plotter.XYs{xys},
			Color: b.FillColor,
		}
		poly.Plot(c, plt)
		
		// Bottom line
		xysBo := plotter.XYs{xys[0], xys[3]}
		lBo := plotter.Line{
			XYs:       xysBo,
			LineStyle: b.LineStyle,
		}
		lBo.Plot(c, plt)
		
		// Upper line
		xysUp := plotter.XYs{xys[1], xys[2]}
		lUp := plotter.Line{
			XYs:       xysUp,
			LineStyle: b.LineStyle,
		}
		lUp.Plot(c, plt)
	}
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (b *BinnedErrBand) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = b.Bins[0][0], b.Bins[len(b.Bins)-1][1]
	ymin, ymax = math.Inf(+1), math.Inf(-1)
	for i, y := range b.Ys {
		ydo, yup := b.YErrs.YError(i)
		ymin = math.Min(ymin, y-ydo)
		ymax = math.Max(ymax, y+yup)
	}
	return xmin, xmax, ymin, ymax
}

var (
	_ plot.Plotter    = (*BinnedErrBand)(nil)
	_ plot.DataRanger = (*BinnedErrBand)(nil)
)
