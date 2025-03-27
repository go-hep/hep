// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"math/rand/v2"
	"testing"
)

func TestCurve2D(t *testing.T) {
	checkPlot(ExampleCurveND_plane, t, "2d-plane-plot.png")
}

func genData2D(n0, n1 int, f func(x, ps []float64) float64, ps []float64, x0min, x0max, x1min, x1max float64) ([][]float64, []float64) {
	var (
		xdata  = make([][]float64, n0*n1)
		ydata  = make([]float64, n0*n1)
		rnd    = rand.New(rand.NewPCG(1234, 1234))
		x0step = (x0max - x0min) / float64(n0)
		x1step = (x1max - x1min) / float64(n1)
		p      = make([]float64, len(ps))
	)
	for i := range n0 {
		for j := range n1 {
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
