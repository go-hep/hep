// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"flag"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/cmpimg"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgtex"
)

var (
	regen = flag.Bool("regen", false, "regenerate reference files")
)

func checkPlot(t *testing.T, ref string) {
	fname := strings.Replace(ref, "_golden", "", 1)

	if *regen {
		got, _ := ioutil.ReadFile(fname)
		ioutil.WriteFile(ref, got, 0644)
	}

	want, err := ioutil.ReadFile(ref)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	ext := filepath.Ext(ref)[1:]
	if ok, err := cmpimg.Equal(ext, got, want); !ok || err != nil {
		if err != nil {
			t.Fatalf("error: comparing %q with reference file: %v\n", fname, err)
		} else {
			t.Fatalf("error: %q differ with reference file\n", fname)
		}
	}
	os.Remove(fname)
}

// An example of a plot + sub-plot
func Example_subplot() {
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

	// normalize histo
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p1 := hplot.New()
	p1.Title.Text = "Histogram"
	p1.X.Label.Text = "X"
	p1.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	p1.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p1.Add(norm)

	// draw a grid
	p1.Add(hplot.NewGrid())

	// make a second plot which will be diplayed in the upper-right
	// of the previous one
	p2 := hplot.New()
	p2.Title.Text = "Sub plot"
	p2.Add(h)
	p2.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgimg.PngCanvas{Canvas: vgimg.New(width, height)}
	dc := draw.New(c)
	p1.Draw(dc)
	sub := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0.70 * width, Y: 0.50 * height},
			Max: vg.Point{X: 1.00 * width, Y: 1.00 * height},
		},
	}
	p2.Draw(sub)

	f, err := os.Create("testdata/sub_plot.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func TestSubPlot(t *testing.T) {
	Example_subplot()
	checkPlot(t, "testdata/sub_plot_golden.png")
}

func Example_diffplot() {

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	hist1 := hbook.NewH1D(20, -4, +4)
	hist2 := hbook.NewH1D(20, -4, +4)

	for i := 0; i < npoints; i++ {
		v1 := dist.Rand()
		v2 := dist.Rand() + 0.5
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	// Make a plot and set its title.
	p1 := hplot.New()
	p1.Title.Text = "Histos"
	p1.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1 := hplot.NewH1D(hist1)
	h1.LineStyle.Color = color.RGBA{R: 255, A: 255}
	h1.FillColor = nil
	p1.Add(h1)

	h2 := hplot.NewH1D(hist2)
	h2.LineStyle.Color = color.RGBA{B: 255, A: 255}
	h2.FillColor = nil
	p1.Add(h2)

	// hide X-axis labels
	p1.X.Tick.Marker = hplot.NoTicks{}

	p1.Add(hplot.NewGrid())

	hist3 := hbook.NewH1D(20, -4, +4)
	for i := 0; i < hist3.Len(); i++ {
		v1 := hist1.Value(i)
		v2 := hist2.Value(i)
		x1, _ := hist1.XY(i)
		hist3.Fill(x1, v1-v2)
	}

	hdiff := hplot.NewH1D(hist3)

	p2 := hplot.New()
	p2.X.Label.Text = "X"
	p2.Y.Label.Text = "Delta-Y"
	p2.Add(hdiff)
	p2.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgimg.PngCanvas{Canvas: vgimg.New(width, height)}
	dc := draw.New(c)
	top := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0.3 * height},
			Max: vg.Point{X: width, Y: height},
		},
	}
	p1.Draw(top)

	bottom := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: width, Y: 0.3 * height},
		},
	}
	p2.Draw(bottom)

	f, err := os.Create("testdata/diff_plot.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestDiffPlot(t *testing.T) {
	Example_diffplot()
	checkPlot(t, "testdata/diff_plot_golden.png")
}

func Example_latexplot() {

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = `Gaussian distribution: $f(x) = \frac{e^{-(x - \mu)^{2}/(2\sigma^{2}) }} {\sigma\sqrt{2\pi}}$`
	p.Y.Label.Text = `$f(x)$`
	p.X.Label.Text = `$x$`

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	h.LineStyle.Color = color.RGBA{R: 255, A: 255}
	h.FillColor = nil
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	p.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgtex.NewDocument(width, height)
	p.Draw(draw.New(c))
	f, err := os.Create("testdata/latex_plot.tex")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestLatexPlot(t *testing.T) {
	Example_latexplot()
	ref, err := ioutil.ReadFile("testdata/latex_plot_golden.tex")
	if err != nil {
		t.Fatal(err)
	}
	chk, err := ioutil.ReadFile("testdata/latex_plot.tex")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ref, chk) {
		t.Fatal("files testdata/latex_plot{,_golden}.tex differ\n")
	}
	os.Remove("testdata/latex_plot.tex")
}
