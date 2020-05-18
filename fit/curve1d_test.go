// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestCurve1D(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleCurve1D_gaussian, t, "gauss-plot.png")
	checkPlot(cmpimg.CheckPlot)(ExampleCurve1D_exponential, t, "exp-plot.png")
	checkPlot(cmpimg.CheckPlot)(ExampleCurve1D_poly, t, "poly-plot.png")
	checkPlot(cmpimg.CheckPlot)(ExampleCurve1D_powerlaw, t, "powerlaw-plot.png")
}

func TestCurve1DGaussianDefaultOpt(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		var (
			cst   = 3.0
			mean  = 30.0
			sigma = 20.0
			want  = []float64{cst, mean, sigma}
		)

		xdata, ydata, err := readXY("testdata/gauss-data.txt")
		if err != nil {
			t.Fatal(err)
		}

		gauss := func(x, cst, mu, sigma float64) float64 {
			v := (x - mu)
			return cst * math.Exp(-v*v/sigma)
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
			nil, nil,
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
			p := hplot.New()
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

			err := p.Save(20*vg.Centimeter, -1, "testdata/gauss-plot-default-opt.png")
			if err != nil {
				t.Fatal(err)
			}
		}
	}, t, "gauss-plot-default-opt.png")
}

func genXY(n int, f func(x float64, ps []float64) float64, ps []float64, xmin, xmax float64) ([]float64, []float64) {
	xdata := make([]float64, n)
	ydata := make([]float64, n)
	rnd := rand.New(rand.NewSource(1234))
	xstep := (xmax - xmin) / float64(n)
	p := make([]float64, len(ps))
	for i := 0; i < n; i++ {
		x := xmin + xstep*float64(i)
		for j := range p {
			v := rnd.NormFloat64()
			p[j] = ps[j] + v*0.2
		}
		xdata[i] = x
		ydata[i] = f(x, p)
	}
	return xdata, ydata
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

func readXYerr(fname string) (xs, ys, yerrs []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, yerrs, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		xs = append(xs, x)

		y, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		ys = append(ys, y)

		yerr, err := strconv.ParseFloat(toks[2], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		yerrs = append(yerrs, yerr)
	}

	return
}
