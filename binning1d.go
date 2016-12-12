// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// binning1D is a 1-dim binning of the x-axis.
type binning1D struct {
	bins     []Bin1D
	dist     dist1D
	outflows [2]dist1D
	xrange   Range
	xstep    float64
}

func newBinning1D(n int, xmin, xmax float64) binning1D {
	if xmin >= xmax {
		panic("hbook: invalid X-axis limits")
	}
	if n <= 0 {
		panic("hbook: X-axis with zero bins")
	}
	bng := binning1D{
		bins:   make([]Bin1D, n),
		xrange: Range{Min: xmin, Max: xmax},
	}
	bng.xstep = float64(n) / bng.xrange.Width()
	for i := range bng.bins {
		bin := &bng.bins[i]
		bin.xrange.Min = xmin + float64(i)/bng.xstep
		bin.xrange.Max = bin.xrange.Min + 1.0/bng.xstep
	}

	return bng
}

func (bng *binning1D) entries() int64 {
	return bng.dist.Entries()
}

func (bng *binning1D) effEntries() float64 {
	return bng.dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *binning1D) xMin() float64 {
	return bng.xrange.Min
}

// xMax returns the high edge of the X-axis
func (bng *binning1D) xMax() float64 {
	return bng.xrange.Max
}

func (bng *binning1D) fill(x, w float64) {
	idx := bng.coordToIndex(x)
	bng.dist.fill(x, w)
	if idx < 0 {
		bng.outflows[-idx-1].fill(x, w)
		return
	}
	bng.bins[idx].fill(x, w)
}

// coordToIndex returns the bin index corresponding to the coordinate x.
func (bng *binning1D) coordToIndex(x float64) int {
	switch {
	default:
		i := int((x - bng.xrange.Min) * bng.xstep)
		return i
	case x < bng.xrange.Min:
		return UnderflowBin
	case x >= bng.xrange.Max:
		return OverflowBin
	}
}

func (bng *binning1D) scaleW(f float64) {
	bng.dist.scaleW(f)
	bng.outflows[0].scaleW(f)
	bng.outflows[1].scaleW(f)
	for i := range bng.bins {
		bin := &bng.bins[i]
		bin.scaleW(f)
	}
}

// Bins returns the slice of bins for this binning.
func (bng *binning1D) Bins() []Bin1D {
	return bng.bins
}
