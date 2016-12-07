// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "math"

// Indices for the under- and over-flow bins.
const (
	UnderflowBin = -2
	OverflowBin  = -1
)

// AxisKind describes the kind of a given axis (fixed-binning, or variable-size binning)
type AxisKind int

// Enumeration of the known axis kinds.
const (
	FixedBinning AxisKind = iota
	VariableBinning
)

// Axis describes an axis (1D, 2D, ...)
type Axis interface {
	Kind() AxisKind
	LowerEdge() float64
	UpperEdge() float64
	Bins() int
	BinLowerEdge(idx int) float64
	BinUpperEdge(idx int) float64
	BinWidth(idx int) float64
	CoordToIndex(coord float64) int
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
