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
func Wrap(p Drawer, opts ...DrawOption) Drawer {
	wp := &wplot{
		p:     p,
		latex: htex.NoopHandler{},
		DPI:   float64(vgimg.DefaultDPI),
	}
	for _, opt := range opts {
		opt(wp)
	}
	return wp
}

type DrawOption func(p *wplot)

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
	return func(p *wplot) {
		p.border = b
	}
}

// WithLatexHandler allows to enable the automatic generation of PDFs from .tex files.
// To enable the automatic generation of PDFs, use DefaultHandler:
//  WithLatexHandler(htex.DefaultHandler)
func WithLatexHandler(h htex.Handler) DrawOption {
	return func(p *wplot) {
		p.latex = h
	}
}

// WithDPI allows to modify the default DPI of a plot.
func WithDPI(dpi float64) DrawOption {
	return func(p *wplot) {
		p.DPI = dpi
	}
}

type wplot struct {
	p      Drawer
	border Border
	latex  htex.Handler
	DPI    float64
}

func (p *wplot) Draw(dc draw.Canvas) {
	vgtexBorder(dc)

	dc = draw.Crop(dc,
		p.border.Left, -p.border.Right,
		p.border.Bottom, -p.border.Top,
	)

	p.p.Draw(dc)
}

var (
	_ Drawer = (*wplot)(nil)
)
