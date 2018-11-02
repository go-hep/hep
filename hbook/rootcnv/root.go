// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootcnv provides tools to convert ROOT histograms and graphs to go-hep/hbook ones.
package rootcnv

import (
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
)

// H1D creates a new H1D from a TH1x.
func H1D(h1 rhist.H1) (*hbook.H1D, error) {
	raw, err := h1.(yodacnv.Marshaler).MarshalYODA()
	if err != nil {
		return nil, err
	}
	var h hbook.H1D
	err = h.UnmarshalYODA(raw)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// H2D creates a new H2D from a TH2x.
func H2D(h2 rhist.H2) (*hbook.H2D, error) {
	raw, err := h2.(yodacnv.Marshaler).MarshalYODA()
	if err != nil {
		return nil, err
	}
	var h hbook.H2D
	err = h.UnmarshalYODA(raw)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// S2D creates a new S2D from a TGraph, TGraphErrors or TGraphAsymmErrors.
func S2D(g rhist.Graph) (*hbook.S2D, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}

	if g, ok := g.(rhist.GraphErrors); ok {
		for i := range pts {
			xlo, xhi := g.XError(i)
			ylo, yhi := g.YError(i)
			pt := &pts[i]
			pt.ErrX = hbook.Range{Min: xlo, Max: xhi}
			pt.ErrY = hbook.Range{Min: ylo, Max: yhi}
		}
	}
	s2d := hbook.NewS2D(pts...)
	s2d.Annotation()["name"] = g.Name()
	s2d.Annotation()["title"] = g.Title()
	return s2d, nil
}

// FromH1D creates a new ROOT TH1D from a 1-dim hbook histogram.
func FromH1D(h1 *hbook.H1D) *rhist.H1D {
	return rhist.NewH1DFrom(h1)
}

// FromH2D creates a new ROOT TH2D from a 2-dim hbook histogram.
func FromH2D(h2 *hbook.H2D) *rhist.H2D {
	return rhist.NewH2DFrom(h2)
}

// FromS2D creates a new ROOT TGraphAsymmErrors from 2-dim hbook data points.
func FromS2D(s2 *hbook.S2D) rhist.GraphErrors {
	return rhist.NewGraphAsymmErrorsFrom(s2)
}
