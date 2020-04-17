// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"go-hep.org/x/hep/hplot/htex"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// Wrap wraps a plot with plot-level drawing options and
// returns a value implementing the Drawer interface.
func Wrap(p Drawer, opts ...DrawOption) *P {
	plt := &P{
		Plot:  p,
		Latex: htex.NoopHandler{},
		DPI:   float64(vgimg.DefaultDPI),
	}
	for _, opt := range opts {
		opt(plt)
	}
	return plt
}

type DrawOption func(p *P)

// Border specifies the borders' sizes, the space between the
// end of the plot image (PDF, PNG, ...) and the actual plot.
type Border struct {
	Left   vg.Length
	Right  vg.Length
	Bottom vg.Length
	Top    vg.Length
}

// WithBorder allows to specify the borders' sizes, the space between the
// end of the plot image (PDF, PNG, ...) and the actual plot.
func WithBorder(b Border) DrawOption {
	return func(p *P) {
		p.Border = b
	}
}

// WithLatexHandler allows to enable the automatic generation of PDFs from .tex files.
// To enable the automatic generation of PDFs, use DefaultHandler:
//  WithLatexHandler(htex.DefaultHandler)
func WithLatexHandler(h htex.Handler) DrawOption {
	return func(p *P) {
		p.Latex = h
	}
}

// WithDPI allows to modify the default DPI of a plot.
func WithDPI(dpi float64) DrawOption {
	return func(p *P) {
		p.DPI = dpi
	}
}

// P is a plot wrapper, holding plot-level customizations.
type P struct {
	// Plot is a gonum/plot.Plot like value.
	Plot Drawer

	// Border specifies the borders' sizes, the space between the
	// end of the plot image (PDF, PNG, ...) and the actual plot.
	Border Border

	// Latex handles the generation of PDFs from .tex files.
	// The default is to use htex.NoopHandler (a no-op).
	// To enable the automatic generation of PDFs, use DefaultHandler:
	//  p := hplot.Wrap(plt)
	//  p.Latex = htex.DefaultHandler
	Latex htex.Handler

	// DPI is the dot-per-inch for PNG,JPEG,... plots.
	DPI float64
}

func (p *P) Draw(dc draw.Canvas) {
	vgtexBorder(dc)

	dc = draw.Crop(dc,
		p.Border.Left, -p.Border.Right,
		p.Border.Bottom, -p.Border.Top,
	)

	p.Plot.Draw(dc)
}

var (
	_ Drawer = (*P)(nil)
)
