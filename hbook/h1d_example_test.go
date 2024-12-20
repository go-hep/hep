// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"math"
	"math/rand/v2"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewPCG(0, 0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for range npoints {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	// fill h with a slice of values and their weights
	h.FillN([]float64{1, 2, 3}, []float64{1, 1, 1})
	h.FillN([]float64{1, 2, 3}, nil) // all weights are 1.

	fmt.Printf("mean:    %.12f\n", h.XMean())
	fmt.Printf("rms:     %.12f\n", h.XRMS())
	fmt.Printf("std-dev: %.12f\n", h.XStdDev())
	fmt.Printf("std-err: %.12f\n", h.XStdErr())

	// Output:
	// mean:    0.002104228518
	// rms:     1.000617135827
	// std-dev: 1.000664927794
	// std-err: 0.010003648633
}

func ExampleAddH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsum := hbook.AddH1D(h1, h2)
	fmt.Printf("Under: %.1f +/- %.1f\n", hsum.Binning.Outflows[0].SumW(), math.Sqrt(hsum.Binning.Outflows[0].SumW2()))
	for i := range hsum.Len() {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, hsum.Binning.Bins[i].SumW(), math.Sqrt(hsum.Binning.Bins[i].SumW2()))
	}
	fmt.Printf("Over : %.1f +/- %.1f\n", hsum.Binning.Outflows[1].SumW(), math.Sqrt(hsum.Binning.Outflows[1].SumW2()))

	// Output:
	// Under: 1.7 +/- 1.2
	// Bin 0: 4.7 +/- 2.4
	// Bin 1: 1.8 +/- 1.3
	// Bin 2: 2.7 +/- 2.1
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: 2.3 +/- 1.6
	// Bin 5: 3.0 +/- 2.2
	// Over : 1.5 +/- 1.1
}

func ExampleAddScaledH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsum := hbook.AddScaledH1D(h1, 10, h2)
	fmt.Printf("Under: %.1f +/- %.1f\n", hsum.Binning.Outflows[0].SumW(), math.Sqrt(hsum.Binning.Outflows[0].SumW2()))
	for i := range hsum.Len() {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, hsum.Binning.Bins[i].SumW(), math.Sqrt(hsum.Binning.Bins[i].SumW2()))
	}
	fmt.Printf("Over : %.1f +/- %.1f\n", hsum.Binning.Outflows[1].SumW(), math.Sqrt(hsum.Binning.Outflows[1].SumW2()))

	// Output:
	// Under: 8.0 +/- 7.1
	// Bin 0: 24.5 +/- 15.7
	// Bin 1: 9.0 +/- 8.1
	// Bin 2: 9.0 +/- 7.3
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: 14.0 +/- 13.0
	// Bin 5: 21.0 +/- 20.0
	// Over : 10.5 +/- 10.0
}

func ExampleSubH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsub := hbook.SubH1D(h1, h2)
	under := hsub.Binning.Outflows[0]
	fmt.Printf("Under: %.1f +/- %.1f\n", under.SumW(), math.Sqrt(under.SumW2()))
	for i, bin := range hsub.Binning.Bins {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, bin.SumW(), math.Sqrt(bin.SumW2()))
	}
	over := hsub.Binning.Outflows[1]
	fmt.Printf("Over : %.1f +/- %.1f\n", over.SumW(), math.Sqrt(over.SumW2()))

	// Output:
	// Under: 0.3 +/- 1.2
	// Bin 0: 0.3 +/- 2.4
	// Bin 1: 0.2 +/- 1.3
	// Bin 2: 1.3 +/- 2.1
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: -0.3 +/- 1.6
	// Bin 5: -1.0 +/- 2.2
	// Over : -0.5 +/- 1.1
}
