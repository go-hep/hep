// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"image/color"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestH1D(t *testing.T) {
	ExampleH1D_gaussian(t)
}

func ExampleH1D_gaussian(t *testing.T) {
	var (
		mean  = 2.0
		sigma = 4.0
		want  = []float64{4.53720e+02, 1.93218e+00, 3.93188e+00} // from ROOT
	)

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:     mean,
		Sigma:  sigma,
		Source: rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(100, -20, +25)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	gauss := func(x, cst, mu, sigma float64) float64 {
		v := (x - mu) / sigma
		return cst * math.Exp(-0.5*v*v)
	}

	res, err := fit.H1D(
		hist,
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return gauss(x, ps[0], ps[1], ps[2])
			},
			N: len(want),
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		t.Fatal(err)
	}
	if got := res.X; !floats.EqualApprox(got, want, 1e-3) {
		t.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p, err := hplot.New()
		if err != nil {
			t.Fatal(err)
		}
		p.X.Label.Text = "f(x) = cst * exp(-0.5 * ((x-mu)/sigma)^2)"
		p.Y.Label.Text = "y-data"
		p.Y.Min = 0

		h, err := hplot.NewH1D(hist)
		if err != nil {
			t.Fatal(err)
		}
		h.Color = color.RGBA{0, 0, 255, 255}
		p.Add(h)

		f := plotter.NewFunction(func(x float64) float64 {
			return gauss(x, res.X[0], res.X[1], res.X[2])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err = p.Save(20*vg.Centimeter, -1, "testdata/h1d-gauss-plot.png")
		if err != nil {
			t.Fatal(err)
		}
	}
}
