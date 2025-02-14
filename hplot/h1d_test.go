// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestH1D(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}
	checkPlot(cmpimg.CheckPlot)(ExampleH1D, t, "h1d_plot.png")
}

func TestH1DtoPDF(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_toPDF, t, "h1d_plot.pdf")
}

func TestH1DLogScale(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_logScaleY, t, "h1d_logy.png")
}

func TestH1DYErrs(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_withYErrBars, t, "h1d_yerrs.png")
}

func TestH1DYErrsBand(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_withYErrBars_withBand, t, "h1d_yerrs_band.png")
}

func TestH1DAsData(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_withYErrBarsAndData, t, "h1d_glyphs.png")
}

func TestH1DLegendStyle(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	checkPlot(cmpimg.CheckPlot)(ExampleH1D_legendStyle, t, "h1d_legend.png")
}

func TestH1DWithBorders(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skipf("ignore test b/c of darwin+Mac-silicon")
	}

	_ = os.Remove("testdata/h1d_borders.png")
	checkPlot(cmpimg.CheckPlot)(ExampleH1D_withPlotBorders, t, "h1d_borders.png")

	_ = os.Remove("testdata/h1d_borders.png")
	// check that it works with a vg.Canvas-WriterTo.
	checkPlot(cmpimg.CheckPlot)(func() {
		const npoints = 10000

		// Create a normal distribution.
		dist := distuv.Normal{
			Mu:    0,
			Sigma: 1,
			Src:   rand.New(rand.NewPCG(0, 0)),
		}

		// Draw some random values from the standard
		// normal distribution.
		hist := hbook.NewH1D(20, -4, +4)
		for range npoints {
			v := dist.Rand()
			hist.Fill(v, 1)
		}

		// normalize histogram
		area := 0.0
		for _, bin := range hist.Binning.Bins {
			area += bin.SumW() * bin.XWidth()
		}
		hist.Scale(1 / area)

		// Make a plot and set its title.
		p := hplot.New()
		p.Title.Text = "Histogram"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		// Create a histogram of our values drawn
		// from the standard normal.
		h := hplot.NewH1D(hist)
		h.Infos.Style = hplot.HInfoSummary
		p.Add(h)

		// The normal distribution function
		norm := hplot.NewFunction(dist.Prob)
		norm.Color = color.RGBA{R: 255, A: 255}
		norm.Width = vg.Points(2)
		p.Add(norm)

		// draw a grid
		p.Add(hplot.NewGrid())

		fig := hplot.Figure(p, hplot.WithBorder(hplot.Border{
			Right:  25,
			Left:   20,
			Top:    25,
			Bottom: 20,
		}))

		c := vgimg.NewWith(
			vgimg.UseWH(6*vg.Inch, 6*vg.Inch/math.Phi),
		)
		dc := draw.New(c)
		fig.Draw(dc)

		f, err := os.Create("testdata/h1d_borders.png")
		if err != nil {
			t.Fatalf("could not create output plot file: %+v", err)
		}
		defer f.Close()

		img := vgimg.PngCanvas{Canvas: c}
		_, err = img.WriteTo(f)
		if err != nil {
			t.Fatalf("could not encode canvas to png: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not save file: %+v", err)
		}
	}, t, "h1d_borders.png")
}
