// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func TestVLine(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleVLine, t, "vline.png")
}

func TestHLine(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleHLine, t, "hline.png")
}

func TestHLineOutOfPlot(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		p := hplot.New()

		pts := []hbook.Point2D{
			{X: 1, Y: 1},
			{X: 2, Y: 2},
		}

		s2d := hplot.NewS2D(hbook.NewS2D(pts...))
		s2d.Color = color.RGBA{R: 255, A: 255}
		s2d.GlyphStyle.Radius = vg.Points(4)

		p.Add(
			s2d, hplot.NewGrid(),
			hplot.HLine(0, nil, color.Gray16{}),
			hplot.HLine(3, color.Gray16{}, nil),
		)

		err := p.Save(5*vg.Centimeter, -1, "testdata/hline_out_of_plot.png")
		if err != nil {
			t.Fatalf("could not save plot: %+v", err)
		}
	}, t, "hline_out_of_plot.png")
}

func TestVLineOutOfPlot(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		p := hplot.New()

		pts := []hbook.Point2D{
			{X: 1, Y: 1},
			{X: 2, Y: 2},
		}

		s2d := hplot.NewS2D(hbook.NewS2D(pts...))
		s2d.Color = color.RGBA{R: 255, A: 255}
		s2d.GlyphStyle.Radius = vg.Points(4)

		p.Add(
			s2d, hplot.NewGrid(),
			hplot.VLine(3, nil, color.Gray16{}),
			hplot.VLine(0, color.Gray16{}, nil),
		)

		err := p.Save(5*vg.Centimeter, -1, "testdata/vline_out_of_plot.png")
		if err != nil {
			t.Fatalf("could not save plot: %+v", err)
		}
	}, t, "vline_out_of_plot.png")
}

func TestHVLineThumbnail(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		p := hplot.New()
		p.Title.Text = "hvlines"
		p.X.Min = 0
		p.X.Max = 10
		p.Y.Min = 0
		p.Y.Max = 10

		var (
			left   = color.Transparent
			right  = color.Transparent
			top    = color.Transparent
			bottom = color.Transparent
		)

		l1 := hplot.VLine(2.5, left, nil)
		l2 := hplot.VLine(5, nil, nil)
		l3 := hplot.VLine(7.5, nil, right)
		l4 := hplot.HLine(2.5, nil, bottom)
		l5 := hplot.HLine(5, nil, nil)
		l6 := hplot.HLine(7.5, top, nil)

		p.Add(l1, l2, l3, l4, l5, l6)
		p.Legend.Add("l1", l1)
		p.Legend.Add("l2", l2)
		p.Legend.Add("l3", l3)
		p.Legend.Add("l4", l4)
		p.Legend.Add("l5", l5)
		p.Legend.Add("l6", l6)

		err := p.Save(-1, -1, "testdata/hvline.png")
		if err != nil {
			t.Fatalf("could not save hvline: %+v", err)
		}
	}, t, "hvline.png")
}
