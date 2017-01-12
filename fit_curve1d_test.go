// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	"image/color"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-hep/fit"
	"github.com/go-hep/hplot"
	"github.com/gonum/floats"
	"github.com/gonum/optimize"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func TestCurve1D(t *testing.T) {
	ExampleCurve1D_gaussian(t)
}

func ExampleCurve1D_gaussian(t *testing.T) {
	const (
		height = 3.0
		mean   = 30.0
		sigma  = 20.0
		ndf    = 2.0
	)

	xdata, ydata, err := readXY("testdata/gauss-data.txt")

	gauss := func(x, mu, sigma, height float64) float64 {
		v := (x - mu)
		return height * math.Exp(-v*v/sigma)
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return gauss(x, ps[0], ps[1], ps[2])
			},
			X:  xdata,
			Y:  ydata,
			Ps: []float64{10, 10, 10},
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		t.Fatal(err)
	}
	if got, want := res.X, []float64{mean, sigma, height}; !floats.EqualApprox(got, want, 1e-3) {
		t.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p, err := hplot.New()
		if err != nil {
			t.Fatal(err)
		}
		p.X.Label.Text = "Gauss"
		p.Y.Label.Text = "y-data"

		s := hplot.NewS2D(hplot.ZipXY(xdata, ydata))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return gauss(x, res.X[0], res.X[1], res.X[2])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err = p.Save(20*vg.Centimeter, -1, "testdata/gauss-plot.png")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func readXY(fname string) (xs, ys []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, err
		}
		xs = append(xs, x)

		y, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, err
		}
		ys = append(ys, y)
	}

	return
}
