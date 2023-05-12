// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !cross_compile

package main

import (
	"image/color"
	"math/rand"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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
			unit.Dp(float32(w.Dots(dpi))),
			unit.Dp(float32(h.Dots(dpi))),
		),
	)
	defer os.Exit(0)

	keys := key.Set(strings.Join(
		[]string{key.NameEscape, "Q", " ", key.NameReturn, key.NameEnter},
		"|",
	))

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			var (
				ops op.Ops
				gtx = layout.NewContext(&ops, e)
			)
			// register a global key listener for the escape key wrapping our entire UI.
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			key.InputOp{
				Tag:  win,
				Keys: keys,
			}.Add(gtx.Ops)

			for _, e := range gtx.Events(win) {
				switch e := e.(type) {
				case key.Event:
					switch e.Name {
					case "Q", key.NameEscape:
						return
					case " ", key.NameReturn, key.NameEnter:
						if e.State == key.Press {
							win.Invalidate()
						}
					}
				}
			}
			area.Pop()

			c := vggio.New(gtx, w, h)
			p := newPlot()
			p.Draw(draw.New(c))
			e.Frame(c.Paint())

		case system.DestroyEvent:
			return
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
