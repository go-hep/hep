// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "sort"

// indices for the 2D-binning overflows
const (
	BngNW int = 1 + iota
	BngN
	BngNE
	BngE
	BngSE
	BngS
	BngSW
	BngW
)

type Binning2D struct {
	Bins     []Bin2D
	Dist     Dist2D
	Outflows [8]Dist2D
	XRange   Range
	YRange   Range
	Nx       int
	Ny       int
	XEdges   []Bin1D
	YEdges   []Bin1D
}

func newBinning2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) Binning2D {
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
	bng := Binning2D{
		Bins:   make([]Bin2D, nx*ny),
		XRange: Range{Min: xlow, Max: xhigh},
		YRange: Range{Min: ylow, Max: yhigh},
		Nx:     nx,
		Ny:     ny,
		XEdges: make([]Bin1D, nx),
		YEdges: make([]Bin1D, ny),
	}
	xwidth := bng.XRange.Width() / float64(bng.Nx)
	ywidth := bng.YRange.Width() / float64(bng.Ny)
	xmin := bng.XRange.Min
	ymin := bng.YRange.Min
	for ix := range bng.XEdges {
		xbin := &bng.XEdges[ix]
		xbin.Range.Min = xmin + float64(ix)*xwidth
		xbin.Range.Max = xmin + float64(ix+1)*xwidth
		for iy := range bng.YEdges {
			ybin := &bng.YEdges[iy]
			ybin.Range.Min = ymin + float64(iy)*ywidth
			ybin.Range.Max = ymin + float64(iy+1)*ywidth
			i := iy*nx + ix
			bin := &bng.Bins[i]
			bin.XRange.Min = xbin.Range.Min
			bin.XRange.Max = xbin.Range.Max
			bin.YRange.Min = ybin.Range.Min
			bin.YRange.Max = ybin.Range.Max
		}
	}
	return bng
}

func newBinning2DFromEdges(xedges, yedges []float64) Binning2D {
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
	bng := Binning2D{
		Bins:   make([]Bin2D, nx*ny),
		XRange: Range{Min: xlow, Max: xhigh},
		YRange: Range{Min: ylow, Max: yhigh},
		Nx:     nx,
		Ny:     ny,
		XEdges: make([]Bin1D, nx),
		YEdges: make([]Bin1D, ny),
	}
	for ix, xmin := range xedges[:nx] {
		xmax := xedges[ix+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bng.XEdges[ix].Range.Min = xmin
		bng.XEdges[ix].Range.Max = xmax
		for iy, ymin := range yedges[:ny] {
			ymax := yedges[iy+1]
			if ymin == ymax {
				panic(errDupEdgesYAxis)
			}
			i := iy*nx + ix
			bin := &bng.Bins[i]
			bin.XRange.Min = xmin
			bin.XRange.Max = xmax
			bin.YRange.Min = ymin
			bin.YRange.Max = ymax
		}
	}
	for iy, ymin := range yedges[:ny] {
		ymax := yedges[iy+1]
		bng.YEdges[iy].Range.Min = ymin
		bng.YEdges[iy].Range.Max = ymax
	}
	return bng
}

func (bng *Binning2D) entries() int64 {
	return bng.Dist.Entries()
}

func (bng *Binning2D) effEntries() float64 {
	return bng.Dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *Binning2D) xMin() float64 {
	return bng.XRange.Min
}

// xMax returns the high edge of the X-axis
func (bng *Binning2D) xMax() float64 {
	return bng.XRange.Max
}

// yMin returns the low edge of the Y-axis
func (bng *Binning2D) yMin() float64 {
	return bng.YRange.Min
}

// yMax returns the high edge of the Y-axis
func (bng *Binning2D) yMax() float64 {
	return bng.YRange.Max
}

func (bng *Binning2D) fill(x, y, w float64) {
	idx := bng.coordToIndex(x, y)
	bng.Dist.fill(x, y, w)
	if idx == len(bng.Bins) {
		// GAP bin
		return
	}
	if idx < 0 {
		bng.Outflows[-idx-1].fill(x, y, w)
		return
	}
	bng.Bins[idx].fill(x, y, w)
}

func (bng *Binning2D) coordToIndex(x, y float64) int {
	ix := Bin1Ds(bng.XEdges).IndexOf(x)
	iy := Bin1Ds(bng.YEdges).IndexOf(y)

	switch {
	case ix == bng.Nx && iy == bng.Ny: // GAP
		return len(bng.Bins)
	case ix == OverflowBin1D && iy == OverflowBin1D:
		return -BngNE
	case ix == OverflowBin1D && iy == UnderflowBin1D:
		return -BngSE
	case ix == UnderflowBin1D && iy == UnderflowBin1D:
		return -BngSW
	case ix == UnderflowBin1D && iy == OverflowBin1D:
		return -BngNW
	case ix == OverflowBin1D:
		return -BngE
	case ix == UnderflowBin1D:
		return -BngW
	case iy == OverflowBin1D:
		return -BngN
	case iy == UnderflowBin1D:
		return -BngS
	}
	return iy*bng.Nx + ix
}
