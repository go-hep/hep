// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

//go:generate go get github.com/campoy/embedmd
//go:generate embedmd -w README.md

package hplot // import "go-hep.org/x/hep/hplot"

import (
	"bytes"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Show displays the plot according to format, returning the raw bytes and
// an error, if any.
//
// If format is the empty string, then "png" is selected.
// The list of accepted format strings is the same one than from
// the gonum.org/v1/plot/vg/draw.NewFormattedCanvas function.
func Show(p Drawer, w, h vg.Length, format string) ([]byte, error) {
	w, h = Dims(w, h)

	if format == "" {
		format = "png"
	}

	c, err := draw.NewFormattedCanvas(w, h, format)
	if err != nil {
		return nil, err
	}

	p.Draw(draw.New(c))
	out := new(bytes.Buffer)
	_, err = c.WriteTo(out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// zip zips together 2 slices and implements the plotter.XYer interface.
type zip struct {
	x []float64
	y []float64
}

// Len implements the plotter.XYer interface
func (z zip) Len() int { return len(z.x) }

// XY implements the plotter.XYer interface
func (z zip) XY(i int) (x, y float64) { return z.x[i], z.y[i] }

// ZipXY zips together 2 slices x and y in such a way to implement the
// plotter.XYer interface.
//
// ZipXY panics if the slices are not of the same length.
func ZipXY(x, y []float64) plotter.XYer {
	if len(x) != len(y) {
		panic("hplot: slices length differ")
	}
	return zip{x: x, y: y}
}
