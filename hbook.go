// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// Bin models 1D, 2D, ... bins.
type Bin interface {
	Rank() int           // Number of dimensions of the bin
	Entries() int64      // Number of entries in the bin
	EffEntries() float64 // Effective number of entries in the bin
	SumW() float64       // sum of weights
	SumW2() float64      // sum of squared weights
}

// Annotation is a bag of attributes that are attached to a histogram.
type Annotation map[string]interface{}

// Histogram is an n-dim histogram (with weighted entries)
type Histogram interface {
	// Annotation returns the annotations attached to the
	// histogram. (e.g. name, title, ...)
	Annotation() Annotation

	// Name returns the name of this histogram
	Name() string

	// Rank returns the number of dimensions of this histogram.
	Rank() int

	// Binning returns the binning of this histogram.
	Binning() Binning

	// Entries returns the number of entries of this histogram.
	Entries() int64
}
