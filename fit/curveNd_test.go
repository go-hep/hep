// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	//"image/color"
	"log"
	//"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestCurve2D(t *testing.T) {
	ExampleCurve2D_plane()
}

func ExampleCurve2D_plane() {
	var (
		m1    = 0.3
		m2    = 0.1
		c     = 0.2
		ps    = []float64{m1, m2, c}
		n0    = 30
		n1    = 30
		x0min = -1.
		x0max = 1.
		x1min = -1.
		x1max = 1.
	)

	plane := func(x, ps []float64) float64 {
		return ps[0]*x[0] + ps[1]*x[1] + ps[2]
	}

	xdata, ydata := genData2D(n0, n1, plane, ps, x0min, x0max, x1min, x1max)

	res, err := fit.CurveND(
		fit.FuncND{
			F: func(x []float64, ps []float64) float64 {
				return plane(x, ps)
			},
			X: xdata,
			Y: ydata,
			N: 3,
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got, want := res.X, []float64{m1, m2, c}; !floats.EqualApprox(got, want, 0.1) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p := hplot.New()
		p.X.Label.Text = "x"
		p.Y.Label.Text = "y"
		p.Y.Min = x1min
		p.Y.Max = x1max
		p.X.Min = x0min
		p.X.Max = x0max

		s := hbook.NewH2D(n0, x0min, x0max, n1, x1min, x1max)

		for i := range xdata {
			s.Fill(xdata[i][0], xdata[i][1], ydata[i])
		}

		p.Add(hplot.NewH2D(s, nil))
		p.Add(plotter.NewGrid())

		err := p.Save(20*vg.Centimeter, -1, "testdata/2DPlane-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genData2D(n0 int, n1 int, f func(x []float64, ps []float64) float64, ps []float64, x0min, x0max float64, x1min, x1max float64) ([][]float64, []float64) {
	xdata := make([][]float64, n0*n1)
	ydata := make([]float64, n0*n1)
	rnd := rand.New(rand.NewSource(1234))
	x0step := (x0max - x0min) / float64(n0)
	x1step := (x1max - x1min) / float64(n1)
	p := make([]float64, len(ps))
	for i := 0; i < n0; i++ {
		for j := 0; j < n1; j++ {
			x := []float64{x0min + x0step*float64(i), x1min + x1step*float64(j)}
			for k := range p {
				v := rnd.NormFloat64()
				p[k] = ps[k] + v*0.05
			}
			xdata[(i%n0)*n0+j] = x
			ydata[(i%n0)*n0+j] = f(x, p)
		}

	}
	return xdata, ydata
}

func readData2D(fname string) (xs [][]float64, ys []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x0, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, err
		}

		x1, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, err
		}

		xs = append(xs, []float64{x0, x1})

		y, err := strconv.ParseFloat(toks[2], 64)
		if err != nil {
			return xs, ys, err
		}
		ys = append(ys, y)
	}

	return
}

func readDataErr2D(fname string) (xs [][]float64, ys, yerrs []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, yerrs, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x0, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}

		x1, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}

		xs = append(xs, []float64{x0, x1})

		y, err := strconv.ParseFloat(toks[2], 64)
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
