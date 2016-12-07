// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// indices for the 2D-axis overflows
const (
	ax2d_NW int = iota
	ax2d_N
	ax2d_NE
	ax2d_E
	ax2d_SE
	ax2d_S
	ax2d_SW
	ax2d_W
)

type axis2D struct {
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

func newAxis2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) axis2D {
	if xlow >= xhigh {
		panic("hbook: invalid X-axis limits")
	}
	if ylow >= yhigh {
		panic("hbook: invalid Y-axis limits")
	}
	if nx == 0 {
		panic("hbook: X-axis with zero bins")
	}
	if ny == 0 {
		panic("hbook: Y-axis with zero bins")
	}
	ax := axis2D{
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

func (axis *axis2D) entries() int64 {
	return axis.dist.Entries()
}

func (axis *axis2D) effEntries() float64 {
	return axis.dist.EffEntries()
}

// minX returns the low edge of the X-axis
func (axis *axis2D) minX() float64 {
	return axis.xrange.Min
}

// maxX returns the high edge of the X-axis
func (axis *axis2D) maxX() float64 {
	return axis.xrange.Max
}

// minY returns the low edge of the Y-axis
func (axis *axis2D) minY() float64 {
	return axis.yrange.Min
}

// maxY returns the high edge of the Y-axis
func (axis *axis2D) maxY() float64 {
	return axis.yrange.Max
}

func (axis *axis2D) fill(x, y, w float64) {
	idx := axis.coordToIndex(x, y)
	axis.dist.fill(x, y, w)
	if idx < 0 {
		axis.outflows[-idx].fill(x, y, w)
		return
	}
	axis.bins[idx].fill(x, y, w)
}

func (axis *axis2D) coordToIndex(x, y float64) int {
	switch {
	case axis.xrange.Min <= x && x < axis.xrange.Max && axis.yrange.Min <= y && y < axis.yrange.Max:
		ix := int((x - axis.xrange.Min) * axis.xstep)
		iy := int((y - axis.yrange.Min) * axis.ystep)
		return iy*axis.nx + ix
	case x >= axis.xrange.Max && axis.yrange.Min <= y && y < axis.yrange.Max:
		return -ax2d_E
	case axis.xrange.Min > x && axis.yrange.Min <= y && y < axis.yrange.Max:
		return -ax2d_W
	case axis.xrange.Min <= x && x < axis.xrange.Max && axis.yrange.Min > y:
		return -ax2d_S
	case axis.xrange.Min <= x && x < axis.xrange.Max && y >= axis.yrange.Max:
		return -ax2d_N
	case axis.xrange.Min > x && y >= axis.yrange.Max:
		return -ax2d_NW
	case x >= axis.xrange.Max && y >= axis.yrange.Max:
		return -ax2d_NE
	case axis.xrange.Min > x && y < axis.yrange.Min:
		return -ax2d_SW
	case x >= axis.xrange.Max && y < axis.yrange.Min:
		return -ax2d_SE
	}
	panic("not reachable")
}

// EvenBinAxis2D is an evenly-binned 2D axis
type EvenBinAxis2D struct {
	x EvenBinAxis
	y EvenBinAxis
}

func NewEvenBinAxis2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) *EvenBinAxis2D {
	return &EvenBinAxis2D{
		x: *NewEvenBinAxis(nx, xlow, xhigh),
		y: *NewEvenBinAxis(ny, ylow, yhigh),
	}
}

func (axis *EvenBinAxis2D) Kind() AxisKind {
	return FixedBinning
}
