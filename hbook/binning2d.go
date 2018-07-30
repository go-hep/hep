// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "sort"

// indices for the 2D-binning overflows
const (
	bngNW int = 1 + iota
	bngN
	bngNE
	bngE
	bngSE
	bngS
	bngSW
	bngW
)

type binning2D struct {
	bins     []Bin2D
	dist     dist2D
	outflows [8]dist2D
	xrange   Range
	yrange   Range
	nx       int
	ny       int
	xedges   []Bin1D
	yedges   []Bin1D
}

func newBinning2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) binning2D {
	if xlow >= xhigh {
		panic(errInvalidXAxis)
	}
	if ylow >= yhigh {
		panic(errInvalidYAxis)
	}
	if nx <= 0 {
		panic(errEmptyXAxis)
	}
	if ny <= 0 {
		panic(errEmptyYAxis)
	}
	bng := binning2D{
		bins:   make([]Bin2D, nx*ny),
		xrange: Range{Min: xlow, Max: xhigh},
		yrange: Range{Min: ylow, Max: yhigh},
		nx:     nx,
		ny:     ny,
		xedges: make([]Bin1D, nx),
		yedges: make([]Bin1D, ny),
	}
	xwidth := bng.xrange.Width() / float64(bng.nx)
	ywidth := bng.yrange.Width() / float64(bng.ny)
	xmin := bng.xrange.Min
	ymin := bng.yrange.Min
	for ix := range bng.xedges {
		xbin := &bng.xedges[ix]
		xbin.xrange.Min = xmin + float64(ix)*xwidth
		xbin.xrange.Max = xmin + float64(ix+1)*xwidth
		for iy := range bng.yedges {
			ybin := &bng.yedges[iy]
			ybin.xrange.Min = ymin + float64(iy)*ywidth
			ybin.xrange.Max = ymin + float64(iy+1)*ywidth
			i := iy*nx + ix
			bin := &bng.bins[i]
			bin.xrange.Min = xbin.xrange.Min
			bin.xrange.Max = xbin.xrange.Max
			bin.yrange.Min = ybin.xrange.Min
			bin.yrange.Max = ybin.xrange.Max
		}
	}
	return bng
}

func newBinning2DFromEdges(xedges, yedges []float64) binning2D {
	if len(xedges) <= 1 {
		panic(errShortXAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(xedges)) {
		panic(errNotSortedXAxis)
	}
	if len(yedges) <= 1 {
		panic(errShortYAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(yedges)) {
		panic(errNotSortedYAxis)
	}
	var (
		nx    = len(xedges) - 1
		ny    = len(yedges) - 1
		xlow  = xedges[0]
		xhigh = xedges[nx]
		ylow  = yedges[0]
		yhigh = yedges[ny]
	)
	bng := binning2D{
		bins:   make([]Bin2D, nx*ny),
		xrange: Range{Min: xlow, Max: xhigh},
		yrange: Range{Min: ylow, Max: yhigh},
		nx:     nx,
		ny:     ny,
		xedges: make([]Bin1D, nx),
		yedges: make([]Bin1D, ny),
	}
	for ix, xmin := range xedges[:nx] {
		xmax := xedges[ix+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bng.xedges[ix].xrange.Min = xmin
		bng.xedges[ix].xrange.Max = xmax
		for iy, ymin := range yedges[:ny] {
			ymax := yedges[iy+1]
			if ymin == ymax {
				panic(errDupEdgesYAxis)
			}
			i := iy*nx + ix
			bin := &bng.bins[i]
			bin.xrange.Min = xmin
			bin.xrange.Max = xmax
			bin.yrange.Min = ymin
			bin.yrange.Max = ymax
		}
	}
	for iy, ymin := range yedges[:ny] {
		ymax := yedges[iy+1]
		bng.yedges[iy].xrange.Min = ymin
		bng.yedges[iy].xrange.Max = ymax
	}
	return bng
}

func (bng *binning2D) entries() int64 {
	return bng.dist.Entries()
}

func (bng *binning2D) effEntries() float64 {
	return bng.dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *binning2D) xMin() float64 {
	return bng.xrange.Min
}

// xMax returns the high edge of the X-axis
func (bng *binning2D) xMax() float64 {
	return bng.xrange.Max
}

// yMin returns the low edge of the Y-axis
func (bng *binning2D) yMin() float64 {
	return bng.yrange.Min
}

// yMax returns the high edge of the Y-axis
func (bng *binning2D) yMax() float64 {
	return bng.yrange.Max
}

func (bng *binning2D) fill(x, y, w float64) {
	idx := bng.coordToIndex(x, y)
	bng.dist.fill(x, y, w)
	if idx == len(bng.bins) {
		// GAP bin
		return
	}
	if idx < 0 {
		bng.outflows[-idx-1].fill(x, y, w)
		return
	}
	bng.bins[idx].fill(x, y, w)
}

func (bng *binning2D) coordToIndex(x, y float64) int {
	ix := Bin1Ds(bng.xedges).IndexOf(x)
	iy := Bin1Ds(bng.yedges).IndexOf(y)

	switch {
	case ix == bng.nx && iy == bng.ny: // GAP
		return len(bng.bins)
	case ix == OverflowBin && iy == OverflowBin:
		return -bngNE
	case ix == OverflowBin && iy == UnderflowBin:
		return -bngSE
	case ix == UnderflowBin && iy == UnderflowBin:
		return -bngSW
	case ix == UnderflowBin && iy == OverflowBin:
		return -bngNW
	case ix == OverflowBin:
		return -bngE
	case ix == UnderflowBin:
		return -bngW
	case iy == OverflowBin:
		return -bngN
	case iy == UnderflowBin:
		return -bngS
	}
	return iy*bng.nx + ix
}

// Bins returns the slice of bins for this binning.
func (bng *binning2D) Bins() []Bin2D {
	return bng.bins
}
