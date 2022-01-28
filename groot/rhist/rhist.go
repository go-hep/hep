// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rhist contains the interfaces and definitions of ROOT types related
// to histograms and graphs.
package rhist // import "go-hep.org/x/hep/groot/rhist"

import (
	"go-hep.org/x/hep/groot/root"
)

// Axis describes a ROOT TAxis.
type Axis interface {
	root.Named

	XMin() float64
	XMax() float64
	NBins() int
	XBins() []float64
	BinCenter(int) float64
	BinLowEdge(int) float64
	BinWidth(int) float64
}

// H1 is a 1-dim ROOT histogram
type H1 interface {
	root.Named

	isH1()

	// Entries returns the number of entries for this histogram.
	Entries() float64
	// SumW returns the total sum of weights
	SumW() float64
	// SumW2 returns the total sum of squares of weights
	SumW2() float64
	// SumWX returns the total sum of weights*x
	SumWX() float64
	// SumWX2 returns the total sum of weights*x*x
	SumWX2() float64
	// SumW2s returns the array of sum of squares of weights
	SumW2s() []float64
}

// H2 is a 2-dim ROOT histogram
type H2 interface {
	root.Named

	isH2()

	// Entries returns the number of entries for this histogram.
	Entries() float64
	// SumW returns the total sum of weights
	SumW() float64
	// SumW2 returns the total sum of squares of weights
	SumW2() float64
	// SumWX returns the total sum of weights*x
	SumWX() float64
	// SumWX2 returns the total sum of weights*x*x
	SumWX2() float64
	// SumW2s returns the array of sum of squares of weights
	SumW2s() []float64
	// SumWY returns the total sum of weights*y
	SumWY() float64
	// SumWY2 returns the total sum of weights*y*y
	SumWY2() float64
	// SumWXY returns the total sum of weights*x*y
	SumWXY() float64
}

// Graph describes a ROOT TGraph
type Graph interface {
	root.Named

	Len() int
	XY(i int) (float64, float64)
}

// GraphErrors describes a ROOT TGraphErrors
type GraphErrors interface {
	Graph
	// XError returns two error values for X data.
	XError(i int) (float64, float64)
	// YError returns two error values for Y data.
	YError(i int) (float64, float64)
}

// F1Composition describes a 1-dim functions composition.
type F1Composition interface {
	root.Object

	isF1Composition() // FIXME(sbinet): have a more useful interface?
	// Eval(xs, ps []float64) float64
}
