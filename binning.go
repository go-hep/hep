// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "math"

// Indices for the under- and over-flow 1-dim bins.
const (
	UnderflowBin = -1
	OverflowBin  = -2
)

// BinningKind describes the kind of binning: fixed-binning, or variable-size binning.
type BinningKind int

// Enumeration of the known binning kinds.
const (
	FixedBinning BinningKind = iota
	VariableBinning
)

// Binning describes the binning of a histogram (1D, 2D, ...)
type Binning interface {
	// Kind returns the binning kind (Fixed,Variable)
	Kind() BinningKind
	// LowerEdge returns the lower edge of the binning.
	LowerEdge() float64
	// UpperEdge returns the upper edge of the binning.
	UpperEdge() float64
	// Bins returns the number of bins in the binning.
	Bins() int
	// BinLowerEdge returns the lower edge of the bin at index i.
	// It panics if i is outside the binning range.
	BinLowerEdge(i int) float64
	// BinUpperEdge returns the upper edge of the bin at index i.
	// It panics if i is outside the binning range.
	BinUpperEdge(i int) float64
	// BinWidth returns the width of the bin at index i.
	BinWidth(idx int) float64
	// CoordToIndex returns the bin index corresponding to the coordinate x.
	CoordToIndex(x float64) int
}

// Range is a 1-dim interval [Min, Max].
type Range struct {
	Min float64
	Max float64
}

// Width returns the size of the range.
func (r Range) Width() float64 {
	return math.Abs(r.Max - r.Min)
}
