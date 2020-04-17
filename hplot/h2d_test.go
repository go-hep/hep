// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestH2D(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleH2D, t, "h2d_plot.png")
}

func TestH2DABCD(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		h := hbook.NewH2D(2, 0, 2, 2, 0, 2)
		h.Fill(0, 0, 1)
		h.Fill(1, 0, 2)
		h.Fill(0, 1, 3)
		h.Fill(1, 1, 4)

		p := hplot.New()
		p.Title.Text = "Hist-2D"
		p.X.Label.Text = "x"
		p.X.Min = -1
		p.X.Max = +3
		p.Y.Label.Text = "y"
		p.Y.Min = -1
		p.Y.Max = +3

		p.Add(hplot.NewH2D(h, nil))
		p.Add(plotter.NewGrid())

		err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot_abcd.png")
		if err != nil {
			t.Fatal(err)
		}
	}, t, "h2d_plot_abcd.png")
}

func TestH2MinMax(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		type pair struct{ i, j int }

		const n = 3
		h2 := hbook.NewH2D(n, 0, n, n, 0, n)
		ws := map[pair]float64{
			{0, 0}: 0,
			{0, 1}: 1,
			{0, 2}: 2,
			{1, 0}: 0,
			{1, 1}: 1,
			{1, 2}: 2,
			{2, 0}: 0,
			{2, 1}: 1,
			{2, 2}: 2,
		}
		for i := 0; i < n; i++ {
			ix := float64(i)
			for j := 0; j < n; j++ {
				iy := float64(j)
				p := pair{i, j}
				v := ws[p]
				h2.Fill(ix, iy, v)
			}
		}

		pl := hplot.New()
		pl.X.Min = 0
		pl.X.Max = 3
		pl.Y.Min = 0
		pl.Y.Max = 3
		hh := hplot.NewH2D(h2, labels{})
		hh.HeatMap.Min = 0
		hh.HeatMap.Max = 3
		pl.Add(hh)

		err := pl.Save(5*vg.Centimeter, 5*vg.Centimeter, "testdata/h2d_plot_minmax.png")
		if err != nil {
			t.Fatal(err)
		}

	}, t, "h2d_plot_minmax.png")
}

type labels struct{}

func (labels) Colors() []color.Color {
	return []color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.Black,
	}
}
