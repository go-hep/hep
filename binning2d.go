// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// indices for the 2D-binning overflows
const (
	bngNW int = iota
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
	xstep    float64
	ystep    float64
}

func newBinning2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) binning2D {
	if xlow >= xhigh {
		panic("hbook: invalid X-axis limits")
	}
	if ylow >= yhigh {
		panic("hbook: invalid Y-axis limits")
	}
	if nx <= 0 {
		panic("hbook: X-axis with zero bins")
	}
	if ny <= 0 {
		panic("hbook: Y-axis with zero bins")
	}
	ax := binning2D{
		bins:   make([]Bin2D, nx*ny),
		xrange: Range{Min: xlow, Max: xhigh},
		yrange: Range{Min: ylow, Max: yhigh},
		nx:     nx,
		ny:     ny,
	}
	ax.xstep = float64(ax.nx) / ax.xrange.Width()
	ax.ystep = float64(ax.ny) / ax.yrange.Width()
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			i := iy*nx + ix
			bin := &ax.bins[i]
			bin.xrange.Min = xlow + float64(ix)/ax.xstep
			bin.xrange.Max = bin.xrange.Min + 1/ax.xstep
			bin.yrange.Min = ylow + float64(iy)/ax.ystep
			bin.yrange.Max = bin.yrange.Min + 1/ax.ystep
		}
	}
	return ax
}

func (bng *binning2D) entries() int64 {
	return bng.dist.Entries()
}

func (bng *binning2D) effEntries() float64 {
	return bng.dist.EffEntries()
}

// minX returns the low edge of the X-axis
func (bng *binning2D) minX() float64 {
	return bng.xrange.Min
}

// maxX returns the high edge of the X-axis
func (bng *binning2D) maxX() float64 {
	return bng.xrange.Max
}

// minY returns the low edge of the Y-axis
func (bng *binning2D) minY() float64 {
	return bng.yrange.Min
}

// maxY returns the high edge of the Y-axis
func (bng *binning2D) maxY() float64 {
	return bng.yrange.Max
}

func (bng *binning2D) fill(x, y, w float64) {
	idx := bng.coordToIndex(x, y)
	bng.dist.fill(x, y, w)
	if idx < 0 {
		bng.outflows[-idx].fill(x, y, w)
		return
	}
	bng.bins[idx].fill(x, y, w)
}

func (bng *binning2D) coordToIndex(x, y float64) int {
	switch {
	case bng.xrange.Min <= x && x < bng.xrange.Max && bng.yrange.Min <= y && y < bng.yrange.Max:
		ix := int((x - bng.xrange.Min) * bng.xstep)
		iy := int((y - bng.yrange.Min) * bng.ystep)
		return iy*bng.nx + ix
	case x >= bng.xrange.Max && bng.yrange.Min <= y && y < bng.yrange.Max:
		return -bngE
	case bng.xrange.Min > x && bng.yrange.Min <= y && y < bng.yrange.Max:
		return -bngW
	case bng.xrange.Min <= x && x < bng.xrange.Max && bng.yrange.Min > y:
		return -bngS
	case bng.xrange.Min <= x && x < bng.xrange.Max && y >= bng.yrange.Max:
		return -bngN
	case bng.xrange.Min > x && y >= bng.yrange.Max:
		return -bngNW
	case x >= bng.xrange.Max && y >= bng.yrange.Max:
		return -bngNE
	case bng.xrange.Min > x && y < bng.yrange.Min:
		return -bngSW
	case x >= bng.xrange.Max && y < bng.yrange.Min:
		return -bngSE
	}
	panic("not reachable")
}
