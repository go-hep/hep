// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/plot/cmpimg"
)

func TestCurve2D(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleCurveND_plane, t, "2d-plane-plot.png")
}

func genData2D(n0, n1 int, f func(x, ps []float64) float64, ps []float64, x0min, x0max, x1min, x1max float64) ([][]float64, []float64) {
	var (
		xdata  = make([][]float64, n0*n1)
		ydata  = make([]float64, n0*n1)
		rnd    = rand.New(rand.NewSource(1234))
		x0step = (x0max - x0min) / float64(n0)
		x1step = (x1max - x1min) / float64(n1)
		p      = make([]float64, len(ps))
	)
	for i := 0; i < n0; i++ {
		for j := 0; j < n1; j++ {
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
