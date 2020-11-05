// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleHStack() {
	h1 := hbook.NewH1D(100, -10, 10)
	h2 := hbook.NewH1D(100, -10, 10)
	h3 := hbook.NewH1D(100, -10, 10)

	const seed = 1234
	fillH1(h1, 10000, -2, 1, seed)
	fillH1(h2, 10000, +3, 3, seed)
	fillH1(h3, 10000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}

	hh1 := hplot.NewH1D(h1)
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black

	hh2 := hplot.NewH1D(h2)
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3)
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black

	hs := []*hplot.H1D{hh1, hh2, hh3}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 3})
	tp.Align = true

	{
		p := tp.Plots[0]
		p.Title.Text = "Histograms"
		p.Y.Label.Text = "Y"
		p.Add(hh1, hh2, hh3, hplot.NewGrid())
		p.Legend.Add("h1", hh1)
		p.Legend.Add("h2", hh2)
		p.Legend.Add("h3", hh3)
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "HStack - stack: OFF"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs)
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(2, 0)
		p.Title.Text = "Hstack - stack: ON"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithLogY(false))
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(15*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

}

func ExampleHStack_withBand() {
	h1 := hbook.NewH1D(50, -8, 12)
	h2 := hbook.NewH1D(50, -8, 12)
	h3 := hbook.NewH1D(50, -8, 12)

	const seed = 1234
	fillH1(h1, 2000, -2, 1, seed)
	fillH1(h2, 2000, +3, 3, seed)
	fillH1(h3, 2000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}

	hh1 := hplot.NewH1D(h1, hplot.WithBand(true))
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black
	hh1.Band.FillColor = color.NRGBA{G: 210, A: 200}

	hh2 := hplot.NewH1D(h2, hplot.WithBand(false))
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3, hplot.WithBand(true))
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black
	hh3.Band.FillColor = color.NRGBA{R: 220, A: 200}

	hs := []*hplot.H1D{hh1, hh2, hh3}

	hh4 := hplot.NewH1D(h1)
	hh4.FillColor = colors[0]
	hh4.LineStyle.Color = color.Black

	hh5 := hplot.NewH1D(h2)
	hh5.FillColor = colors[1]
	hh5.LineStyle.Width = 0

	hh6 := hplot.NewH1D(h3)
	hh6.FillColor = colors[2]
	hh6.LineStyle.Color = color.Black

	hsHistoNoBand := []*hplot.H1D{hh4, hh5, hh6}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 2})
	tp.Align = true

	{
		p := tp.Plot(0, 0)
		p.Title.Text = "Histos With or Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true))
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(0, 1)
		p.Title.Text = "Histos Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true))
		hstack.Stack = hplot.HStackOff
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "Histos With or Without Band, Stack: ON"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true))
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 1)
		p.Title.Text = "Histos Without Band, Stack: ON"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true))
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(25*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack_band.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}

func ExampleHStack_withLogY() {
	h1 := hbook.NewH1D(50, -8, 12)
	h2 := hbook.NewH1D(50, -8, 12)
	h3 := hbook.NewH1D(50, -8, 12)

	const seed = 1234
	fillH1(h1, 2000, -2, 1, seed)
	fillH1(h2, 2000, +3, 3, seed)
	fillH1(h3, 2000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}
	logy := hplot.WithLogY(true)

	hh1 := hplot.NewH1D(h1, hplot.WithBand(true), logy)
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black
	hh1.Band.FillColor = color.NRGBA{G: 210, A: 200}

	hh2 := hplot.NewH1D(h2, hplot.WithBand(false), logy)
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3, hplot.WithBand(true), logy)
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black
	hh3.Band.FillColor = color.NRGBA{R: 220, A: 200}

	hs := []*hplot.H1D{hh1, hh2, hh3}

	hh4 := hplot.NewH1D(h1, logy)
	hh4.FillColor = colors[0]
	hh4.LineStyle.Color = color.Black

	hh5 := hplot.NewH1D(h2, logy)
	hh5.FillColor = colors[1]
	hh5.LineStyle.Width = 0

	hh6 := hplot.NewH1D(h3, logy)
	hh6.FillColor = colors[2]
	hh6.LineStyle.Color = color.Black

	hsHistoNoBand := []*hplot.H1D{hh4, hh5, hh6}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 2})
	tp.Align = true

	{
		p := tp.Plot(0, 0)
		p.Title.Text = "Histos With or Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true), logy)
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(0, 1)
		p.Title.Text = "Histos Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true), logy)
		hstack.Stack = hplot.HStackOff
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "Histos With or Without Band, Stack: ON"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true), logy)
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 1)
		p.Title.Text = "Histos Without Band, Stack: ON"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true), logy)
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(25*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack_logy.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}

func fillH1(h *hbook.H1D, n int, mu, sigma float64, seed uint64) {
	dist := distuv.Normal{
		Mu:    mu,
		Sigma: sigma,
		Src:   rand.New(rand.NewSource(seed)),
	}

	for i := 0; i < n; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
}
