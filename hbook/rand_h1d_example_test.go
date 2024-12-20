// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"path"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func ExampleRand1D() {
	const N = 10000
	var (
		h1  = hbook.NewH1D(100, -10, 10)
		src = rand.New(rand.NewPCG(1234, 1234))
		rnd = distuv.Normal{
			Mu:    0,
			Sigma: 2,
			Src:   src,
		}
	)

	for range N {
		h1.Fill(rnd.Rand(), 1)
	}

	var (
		h2 = hbook.NewH1D(100, -10, 10)
		hr = hbook.NewRand1D(h1, rand.NewPCG(5678, 5678))
	)

	for range N {
		h2.Fill(hr.Rand(), 1)
	}

	fmt.Printf(
		"h1: mean=%+8f std-dev=%+8f +/- %8f\n",
		h1.XMean(), h1.XStdDev(), h1.XStdErr(),
	)
	fmt.Printf(
		"cdf(0)= %+1.1f\ncdf(1)= %+1.1f\n",
		rnd.CDF(0), rnd.CDF(1),
	)
	fmt.Printf(
		"h2: mean=%+8f std-dev=%+8f +/- %8f\n",
		h2.XMean(), h2.XStdDev(), h2.XStdErr(),
	)
	fmt.Printf(
		"cdf(0)= %+1.1f\ncdf(1)= %+1.1f\n",
		hr.CDF(0), hr.CDF(1),
	)

	h1.Scale(1. / h1.Integral(h1.XMin(), h1.XMax()))
	h2.Scale(1. / h2.Integral(h2.XMin(), h2.XMax()))

	{
		rp := hplot.NewRatioPlot()
		rp.Ratio = 0.3

		rp.Top.Title.Text = "Distributions"
		rp.Top.Y.Label.Text = "Y"

		hh1 := hplot.NewH1D(h1)
		hh1.FillColor = color.NRGBA{R: 255, A: 100}

		hh2 := hplot.NewH1D(h2)
		hh2.FillColor = color.NRGBA{B: 255, A: 100}

		rp.Top.Add(
			hh1, hh2,
			hplot.NewGrid(),
		)

		rp.Top.Legend.Add("template", hh1)
		rp.Top.Legend.Add("monte-carlo", hh2)
		rp.Top.Legend.Top = true

		rp.Bottom.X.Label.Text = "X"
		rp.Bottom.Y.Label.Text = "Diff"
		rp.Bottom.Add(
			hplot.NewH1D(hbook.SubH1D(h1, h2)),
			hplot.NewGrid(),
		)

		err := hplot.Save(rp, 15*vg.Centimeter, -1, "testdata/rand_h1d.png")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Output:
	// h1: mean=+0.011317 std-dev=+2.027279 +/- 0.020273
	// cdf(0)= +0.5
	// cdf(1)= +0.7
	// h2: mean=-0.026749 std-dev=+2.022967 +/- 0.020230
	// cdf(0)= +0.5
	// cdf(1)= +0.7
}

func TestRand1DExample(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleRand1D, t, "rand_h1d.png")
}

type chkplotFunc func(ExampleFunc func(), t *testing.T, filenames ...string)

func checkPlot(f chkplotFunc) chkplotFunc {
	return func(ex func(), t *testing.T, filenames ...string) {
		t.Helper()
		f(ex, t, filenames...)
		if t.Failed() {
			return
		}
		for _, fname := range filenames {
			_ = os.Remove(path.Join("testdata", fname))
		}
	}
}

var (
	_ distuv.Rander = (*hbook.Rand1D)(nil)
)
