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
	cmpimg.CheckPlot(ExampleVLine, t, "vline.png")
}

func TestHLine(t *testing.T) {
	cmpimg.CheckPlot(ExampleHLine, t, "hline.png")
}

func TestHLineOutOfPlot(t *testing.T) {
	cmpimg.CheckPlot(func() {
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
	cmpimg.CheckPlot(func() {
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
