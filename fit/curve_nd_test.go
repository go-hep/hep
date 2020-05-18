// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	"fmt"
	"image/color"

	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"go-hep.org/x/hep/fit"
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
		m1         = 0.3
		m2         = 0.1
		c          = 0.2
		ps         = []float64{m1, m2, c}
		n0    uint = 10
		n1    uint = 10
		x0min      = -1.
		x0max      = 1.
		x1min      = -1.
		x1max      = 1.
	)

	plane := func(x, ps []float64) float64 {
		return ps[0]*x[0] + ps[1]*x[1] + ps[2]
	}

	xData, yData := genData2D(n0, n1, plane, ps, x0min, x0max, x1min, x1max)

	res, err := fit.CurveND(
		fit.FuncND{
			F: func(x []float64, ps []float64) float64 {
				return plane(x, ps)
			},
			X: xData,
			Y: yData,
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
		p.X.Label.Text = "x1"
		p.Y.Label.Text = "y"
		p.Y.Min = x1min
		p.Y.Max = x1max
		p.X.Min = x0min
		p.X.Max = x0max

		// slicing for a particular x0 value to plot y as a function of x1, to visualise how well the
		// the fit is working for a given x0.
		var x0Selection uint = 8
		if x0Selection > n0 {
			log.Fatalf("x0 slice, %d, is not in valid range [0 - %d]", x0Selection, n0)
		}
		x0SlicePos := x0min + ((x0max-x0min)/float64(n0))*float64(x0Selection)

		var x1Slice []float64
		var ySlice []float64

		for i := range xData {
			if xData[i][0] == x0SlicePos {
				x1Slice = append(x1Slice, xData[i][1])
				ySlice = append(ySlice, yData[i])
			}
		}

		s := hplot.NewS2D(hplot.ZipXY(x1Slice, ySlice))
		s.Color = color.RGBA{B: 255, A: 255}
		p.Add(s)

		shiftLine := func(x, m, c, mxOtherAxis float64) float64 {
			return m*x + c + mxOtherAxis
		}

		f := plotter.NewFunction(func(x float64) float64 {
			return shiftLine(x, res.X[1], res.X[2], res.X[0]*x0SlicePos)
		})
		f.Color = color.RGBA{R: 255, A: 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())
		p.Title.Text = fmt.Sprintf("Slice of plane at x0 = %.2f", x0SlicePos)
		err := p.Save(20*vg.Centimeter, -1, "testdata/2dplane-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genData2D(n0 uint, n1 uint, f func(x []float64, ps []float64) float64, ps []float64, x0min, x0max float64, x1min, x1max float64) ([][]float64, []float64) {
	xdata := make([][]float64, n0*n1)
	ydata := make([]float64, n0*n1)
	rnd := rand.New(rand.NewSource(1234))
	x0step := (x0max - x0min) / float64(n0)
	x1step := (x1max - x1min) / float64(n1)
	p := make([]float64, len(ps))
	for i := uint(0); i < n0; i++ {
		for j := uint(0); j < n1; j++ {
			x := []float64{x0min + x0step*float64(i), x1min + x1step*float64(j)}
			for k := range p {
				v := rnd.NormFloat64()
				p[k] = ps[k] + v*0.01
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
