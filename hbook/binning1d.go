// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"errors"
	"sort"
)

// Indices for the under- and over-flow 1-dim bins.
const (
	UnderflowBin1D = -1
	OverflowBin1D  = -2
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

// Binning1D is a 1-dim binning of the x-axis.
type Binning1D struct {
	Bins     []Bin1D
	Dist     Dist1D
	Outflows [2]Dist1D
	XRange   Range
}

func newBinning1D(n int, xmin, xmax float64) Binning1D {
	if xmin >= xmax {
		panic(errInvalidXAxis)
	}
	if n <= 0 {
		panic(errEmptyXAxis)
	}
	bng := Binning1D{
		Bins:   make([]Bin1D, n),
		XRange: Range{Min: xmin, Max: xmax},
	}
	width := bng.XRange.Width() / float64(n)
	for i := range bng.Bins {
		bin := &bng.Bins[i]
		bin.Range.Min = xmin + float64(i)*width
		bin.Range.Max = xmin + float64(i+1)*width
	}
	return bng
}

func newBinning1DFromBins(xbins []Range) Binning1D {
	if len(xbins) < 1 {
		panic(errShortXAxis)
	}
	n := len(xbins)
	bng := Binning1D{
		Bins: make([]Bin1D, n),
	}
	for i, xbin := range xbins {
		bin := &bng.Bins[i]
		bin.Range = xbin
	}
	sort.Sort(Bin1Ds(bng.Bins))
	for i := 0; i < len(bng.Bins)-1; i++ {
		b0 := bng.Bins[i]
		b1 := bng.Bins[i+1]
		if b0.Range.Max > b1.Range.Min {
			panic(errOverlapXAxis)
		}
	}
	bng.XRange = Range{Min: bng.Bins[0].XMin(), Max: bng.Bins[n-1].XMax()}
	return bng
}

func newBinning1DFromEdges(edges []float64) Binning1D {
	if len(edges) <= 1 {
		panic(errShortXAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(edges)) {
		panic(errNotSortedXAxis)
	}
	n := len(edges) - 1
	bng := Binning1D{
		Bins:   make([]Bin1D, n),
		XRange: Range{Min: edges[0], Max: edges[n]},
	}
	for i := range bng.Bins {
		bin := &bng.Bins[i]
		xmin := edges[i]
		xmax := edges[i+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bin.Range.Min = xmin
		bin.Range.Max = xmax
	}
	return bng
}

func (bng *Binning1D) clone() Binning1D {
	o := Binning1D{
		Bins: make([]Bin1D, len(bng.Bins)),
		Dist: bng.Dist.clone(),
		Outflows: [2]Dist1D{
			bng.Outflows[0].clone(),
			bng.Outflows[1].clone(),
		},
		XRange: bng.XRange.clone(),
	}

	for i, bin := range bng.Bins {
		o.Bins[i] = bin.clone()
	}

	return o
}

func (bng *Binning1D) entries() int64 {
	return bng.Dist.Entries()
}

func (bng *Binning1D) effEntries() float64 {
	return bng.Dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *Binning1D) xMin() float64 {
	return bng.XRange.Min
}

// xMax returns the high edge of the X-axis
func (bng *Binning1D) xMax() float64 {
	return bng.XRange.Max
}

func (bng *Binning1D) fill(x, w float64) {
	idx := bng.coordToIndex(x)
	bng.Dist.fill(x, w)
	if idx < 0 {
		bng.Outflows[-idx-1].fill(x, w)
		return
	}
	if idx == len(bng.Bins) {
		// gap bin.
		return
	}
	bng.Bins[idx].fill(x, w)
}

// coordToIndex returns the bin index corresponding to the coordinate x.
func (bng *Binning1D) coordToIndex(x float64) int {
	switch {
	case x < bng.XRange.Min:
		return UnderflowBin1D
	case x >= bng.XRange.Max:
		return OverflowBin1D
	}
	return Bin1Ds(bng.Bins).IndexOf(x)
}

func (bng *Binning1D) scaleW(f float64) {
	bng.Dist.scaleW(f)
	bng.Outflows[0].scaleW(f)
	bng.Outflows[1].scaleW(f)
	for i := range bng.Bins {
		bin := &bng.Bins[i]
		bin.scaleW(f)
	}
}

func (bng *Binning1D) Underflow() *Dist1D {
	return &bng.Outflows[0]
}

func (bng *Binning1D) Overflow() *Dist1D {
	return &bng.Outflows[1]
}
