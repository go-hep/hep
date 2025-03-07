// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"log"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExampleH2D() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for range npoints {
		v = dist.Rand(v)
		h.Fill(v[0], v[1], 1)
	}

	p := hplot.New()
	p.Title.Text = "Hist-2D"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	p.Add(hplot.NewH2D(h, nil))
	p.Add(plotter.NewGrid())
	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot.png")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleH2D_withLegend() {
	h2d := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for range npoints {
		v = dist.Rand(v)
		h2d.Fill(v[0], v[1], 1)
	}
	h := hplot.NewH2D(h2d, nil)

	p := hplot.New()
	p.Title.Text = "Hist-2D"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	p.Add(h)
	p.Add(plotter.NewGrid())

	fig := hplot.Figure(p, hplot.WithLegend(h.Legend()))
	err := hplot.Save(fig, 10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot_legend.png")
	if err != nil {
		log.Fatal(err)
	}
}
