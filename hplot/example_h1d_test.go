// Copyright ©2016 The go-hep Authors. All rights reserved.
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
	"gonum.org/v1/plot/plotutil"
)

// An example of making a 1D-histogram.
func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_plot.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// An example of making a 1D-histogram and saving to a PDF
func ExampleH1D_toPDF() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_plot.pdf"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

func ExampleH1D_logScaleY() {
	p := hplot.New()
	p.Title.Text = "Histogram in log-y"
	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Y.Label.Text = "Y"
	p.X.Label.Text = "X"

	h1 := hbook.NewH1D(10, -5, +5)
	for _, v := range []float64{
		-2, -2,
		-1,
		+3, +3, +3, +3,
		+1, +1, +1, +1, +1, +1, +1, +1, +1, +1,
		+1, +1, +1, +1, +1, +1, +1, +1, +1, +1,
	} {
		h1.Fill(v, 1)
	}
	p1 := hplot.NewH1D(h1)
	p1.LogY = true
	p1.FillColor = color.RGBA{255, 0, 0, 255}

	h2 := hbook.NewH1D(10, -5, +5)
	for _, v := range []float64{
		-3, -3, -3,
		+2, +2, +2, +2, +2,
	} {
		h2.Fill(v, 1)
	}
	p2 := hplot.NewH1D(h2,
		hplot.WithYErrBars(true),
		hplot.WithLogY(true),
		hplot.WithGlyphStyle(draw.GlyphStyle{
			Color:  color.Black,
			Radius: 2,
			Shape:  draw.CircleGlyph{},
		}),
	)
	p2.FillColor = color.RGBA{0, 0, 255, 255}

	p.Add(p1, p2, hplot.NewGrid())

	err := p.Save(6*vg.Inch, -1, "testdata/h1d_logy.png")
	if err != nil {
		log.Fatal(err)
	}
}

// An example of making a 1D-histogram with y-error bars.
func ExampleH1D_withYErrBars() {
	const npoints = 100

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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
	h := hplot.NewH1D(hist,
		hplot.WithHInfo(hplot.HInfoSummary),
		hplot.WithYErrBars(true),
	)
	h.YErrs.LineStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_yerrs.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// An example of making a 1D-histogram with y-error bars
// and the associated band.
func ExampleH1D_withYErrBars_withBand() {
	const npoints = 100

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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
	h := hplot.NewH1D(hist,
		hplot.WithHInfo(hplot.HInfoSummary),
		hplot.WithYErrBars(true),
		hplot.WithBand(true),
	)
	h.YErrs.LineStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_yerrs_band.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// An example of making a 1D-histogram with y-error bars
// and no histogram rectangle.
func ExampleH1D_withYErrBarsAndData() {
	const npoints = 100

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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

	p.Legend.Top = true
	p.Legend.Left = true

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist,
		hplot.WithHInfo(hplot.HInfoSummary),
		hplot.WithYErrBars(true),
		hplot.WithGlyphStyle(draw.GlyphStyle{
			Shape:  draw.CrossGlyph{},
			Color:  color.Black,
			Radius: vg.Points(2),
		}),
	)
	h.GlyphStyle.Shape = draw.CircleGlyph{}
	h.YErrs.LineStyle.Color = color.Black
	h.LineStyle.Width = 0 // disable histogram lines
	p.Add(h)
	p.Legend.Add("data", h)
	
	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)
	p.Legend.Add("model", norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_glyphs.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// An example of making a 1D-histogram.
func ExampleH1D_withPlotBorders() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
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

	fig := hplot.Figure(p,
		hplot.WithDPI(96),
		hplot.WithBorder(hplot.Border{
			Right:  25,
			Left:   20,
			Top:    25,
			Bottom: 20,
		}),
	)

	// Save the plot to a PNG file.
	if err := hplot.Save(fig, 6*vg.Inch, -1, "testdata/h1d_borders.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}


// An example showing legend with different style
func ExampleH1D_legendStyle() {

	const npoints = 500
	
	// Create a few normal distributions.
	var hists [3]*hbook.H1D
	for id := range [3]int{0, 1, 2} {
		mu := -2. + float64(id)*2
		sigma := 0.3
			
		dist := distuv.Normal{
			Mu:    mu,
			Sigma: sigma,
			Src:   rand.New(rand.NewSource(uint64(id))),
		}
		
		// Draw some random values from the standard
		// normal distribution.
		hists[id] = hbook.NewH1D(20, mu-5*sigma, mu+5*sigma)
		for i := 0; i < npoints; i++ {
			v := dist.Rand()
			hists[id].Fill(v, 1)
		}
	}

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.X.Max = 10 

	// Legend style tunning
	p.Legend.Top = true
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch
	p.Legend.TextStyle.Font.Size = 12
	p.Legend.Padding = 0.1 * vg.Inch
	
	// Histogram with error band
	hband := hplot.NewH1D(hists[0], hplot.WithBand(true))
	hband.LineStyle.Width = 1.3
	p.Add(hband)
	p.Legend.Add("Band", hband)
	
	// Histogram with fill and line
	hfill := hplot.NewH1D(hists[1])
	hfill.FillColor = color.NRGBA{R: 200, A: 180}
	hfill.LineStyle.Color = color.NRGBA{R: 200, G: 30, A: 255}
	hfill.LineStyle.Width = 1.3
	p.Add(hfill)
	p.Legend.Add("Fill & line", hfill)
	
	// Histogram with line and markers
	hmarker := hplot.NewH1D(hists[2],
		hplot.WithYErrBars(true),
		hplot.WithGlyphStyle(draw.GlyphStyle{
			Color:  color.Black,
			Radius: 2,
			Shape:  draw.CircleGlyph{},
		}),
	)
	hmarker.LineStyle.Color = color.NRGBA{R: 200, G: 30, A: 255}
	hmarker.LineStyle.Width = 1
	hmarker.LineStyle.Dashes = plotutil.Dashes(2)
	p.Add(hmarker)
	p.Legend.Add("marker & line", hmarker)

	// Create a figure
	fig := hplot.Figure(p, hplot.WithDPI(192.))
	fig.Border.Right = 15
	fig.Border.Left = 10
	fig.Border.Top = 10
	fig.Border.Bottom = 10

	// Save the figure to a PNG file.
	if err := hplot.Save(fig, 6*vg.Inch, -1, "testdata/h1d_legend.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}	
}
