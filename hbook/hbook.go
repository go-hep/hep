// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook // import "go-hep.org/x/hep/hbook"

import (
	"math"
)

//go:generate go get github.com/campoy/embedmd
//go:generate embedmd -w README.md

//go:generate brio-gen -p go-hep.org/x/hep/hbook -t Dist0D,Dist1D,Dist2D -o dist_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t Range,Binning1D,binningP1D,Bin1D,BinP1D,Binning2D,Bin2D -o binning_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t Point2D -o points_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t H1D,H2D,P1D,S2D -o hbook_brio.go

// Bin models 1D, 2D, ... bins.
type Bin interface {
	Rank() int           // Number of dimensions of the bin
	Entries() int64      // Number of entries in the bin
	EffEntries() float64 // Effective number of entries in the bin
	SumW() float64       // sum of weights
	SumW2() float64      // sum of squared weights
}

// Range is a 1-dim interval [Min, Max].
type Range struct {
	Min float64
	Max float64
}

func (r Range) clone() Range {
	return r
}

// Width returns the size of the range.
func (r Range) Width() float64 {
	return math.Abs(r.Max - r.Min)
}

// Histogram is an n-dim histogram (with weighted entries)
type Histogram interface {
	// Annotation returns the annotations attached to the
	// histogram. (e.g. name, title, ...)
	Annotation() Annotation

	// Name returns the name of this histogram
	Name() string

	// Rank returns the number of dimensions of this histogram.
	Rank() int

	// Entries returns the number of entries of this histogram.
	Entries() int64
}
