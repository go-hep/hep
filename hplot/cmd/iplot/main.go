// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !cross_compile

package main

import (
	"image/color"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vggio"
)

const (
	NPOINTS = 100000
	dpi     = 96
)

func main() {
	go run()
	app.Main()
}

func run() {
	w, h := hplot.Dims(-1, -1)
	win := app.NewWindow(
		app.Title("iplot"),
		app.Size(
			unit.Px(float32(w.Dots(dpi))),
			unit.Px(float32(h.Dots(dpi))),
		),
	)
	defer win.Close()

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			c := vggio.New(layout.NewContext(new(op.Ops), e), w, h)
			p := newPlot()
			p.Draw(draw.New(c))
			e.Frame(c.Paint())

		case system.DestroyEvent:
			return

		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			case " ", key.NameReturn, key.NameEnter:
				win.Invalidate()
			}
		}
	}
}

func newPlot() *hplot.Plot {
	// Draw some random values from the standard
	// normal distribution.
	hist1 := hbook.NewH1D(100, -5, +5)
	hist2 := hbook.NewH1D(100, -5, +5)
	for i := 0; i < NPOINTS; i++ {
		v1 := rand.NormFloat64() - 1
		v2 := rand.NormFloat64() + 1
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1 := hplot.NewH1D(hist1)
	h1.Infos.Style = hplot.HInfoSummary
	h1.Color = color.Black
	h1.FillColor = nil

	h2 := hplot.NewH1D(hist2)
	h2.Infos.Style = hplot.HInfoNone
	h2.Color = color.RGBA{255, 0, 0, 255}
	h2.FillColor = nil

	p.Add(h1, h2)

	p.Add(plotter.NewGrid())
	return p
}
