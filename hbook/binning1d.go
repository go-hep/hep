// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"errors"
	"sort"
)

// Indices for the under- and over-flow 1-dim bins.
const (
	UnderflowBin = -1
	OverflowBin  = -2
)

var (
	errInvalidXAxis   = errors.New("hbook: invalid X-axis limits")
	errEmptyXAxis     = errors.New("hbook: X-axis with zero bins")
	errShortXAxis     = errors.New("hbook: too few 1-dim X-bins")
	errOverlapXAxis   = errors.New("hbook: invalid X-binning (overlap)")
	errNotSortedXAxis = errors.New("hbook: X-edges slice not sorted")
	errDupEdgesXAxis  = errors.New("hbook: duplicates in X-edge values")

	errInvalidYAxis   = errors.New("hbook: invalid Y-axis limits")
	errEmptyYAxis     = errors.New("hbook: Y-axis with zero bins")
	errShortYAxis     = errors.New("hbook: too few 1-dim Y-bins")
	errOverlapYAxis   = errors.New("hbook: invalid Y-binning (overlap)")
	errNotSortedYAxis = errors.New("hbook: Y-edges slice not sorted")
	errDupEdgesYAxis  = errors.New("hbook: duplicates in Y-edge values")
)

// binning1D is a 1-dim binning of the x-axis.
type binning1D struct {
	bins     []Bin1D
	dist     dist1D
	outflows [2]dist1D
	xrange   Range
}

func newBinning1D(n int, xmin, xmax float64) binning1D {
	if xmin >= xmax {
		panic(errInvalidXAxis)
	}
	if n <= 0 {
		panic(errEmptyXAxis)
	}
	bng := binning1D{
		bins:   make([]Bin1D, n),
		xrange: Range{Min: xmin, Max: xmax},
	}
	width := bng.xrange.Width() / float64(n)
	for i := range bng.bins {
		bin := &bng.bins[i]
		bin.xrange.Min = xmin + float64(i)*width
		bin.xrange.Max = xmin + float64(i+1)*width
	}
	return bng
}

func newBinning1DFromBins(xbins []Range) binning1D {
	if len(xbins) < 1 {
		panic(errShortXAxis)
	}
	n := len(xbins)
	bng := binning1D{
		bins: make([]Bin1D, n),
	}
	for i, xbin := range xbins {
		bin := &bng.bins[i]
		bin.xrange = xbin
	}
	sort.Sort(Bin1Ds(bng.bins))
	for i := 0; i < len(bng.bins)-1; i++ {
		b0 := bng.bins[i]
		b1 := bng.bins[i+1]
		if b0.xrange.Max > b1.xrange.Min {
			panic(errOverlapXAxis)
		}
	}
	bng.xrange = Range{Min: bng.bins[0].XMin(), Max: bng.bins[n-1].XMax()}
	return bng
}

func newBinning1DFromEdges(edges []float64) binning1D {
	if len(edges) <= 1 {
		panic(errShortXAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(edges)) {
		panic(errNotSortedXAxis)
	}
	n := len(edges) - 1
	bng := binning1D{
		bins:   make([]Bin1D, n),
		xrange: Range{Min: edges[0], Max: edges[n]},
	}
	for i := range bng.bins {
		bin := &bng.bins[i]
		xmin := edges[i]
		xmax := edges[i+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bin.xrange.Min = xmin
		bin.xrange.Max = xmax
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
	if idx == len(bng.bins) {
		// gap bin.
		return
	}
	bng.bins[idx].fill(x, w)
}

// coordToIndex returns the bin index corresponding to the coordinate x.
func (bng *binning1D) coordToIndex(x float64) int {
	switch {
	case x < bng.xrange.Min:
		return UnderflowBin
	case x >= bng.xrange.Max:
		return OverflowBin
	}
	return Bin1Ds(bng.bins).IndexOf(x)
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
