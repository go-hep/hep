// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootcnv provides tools to convert ROOT histograms and graphs to go-hep/hbook ones.
package rootcnv

import (
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/rootio"
)

// H1D creates a new H1D from a TH1x.
func H1D(r yodacnv.Marshaler) (*hbook.H1D, error) {
	raw, err := r.MarshalYODA()
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
func H2D(r yodacnv.Marshaler) (*hbook.H2D, error) {
	raw, err := r.MarshalYODA()
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
func S2D(g rootio.Graph) (*hbook.S2D, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}

	if g, ok := g.(rootio.GraphErrors); ok {
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
