// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExampleCurve1D_gaussian() {
	var (
		cst   = 3.0
		mean  = 30.0
		sigma = 20.0
		want  = []float64{cst, mean, sigma}
	)

	xdata, ydata, err := readXY("testdata/gauss-data.txt")
	if err != nil {
		log.Fatal(err)
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
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got := res.X; !floats.EqualApprox(got, want, 1e-3) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
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

		err := p.Save(20*vg.Centimeter, -1, "testdata/gauss-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_exponential() {
	const (
		a   = 0.3
		b   = 0.1
		ndf = 2.0
	)

	xdata, ydata, err := readXY("testdata/exp-data.txt")
	if err != nil {
		log.Fatal(err)
	}

	exp := func(x, a, b float64) float64 {
		return math.Exp(a*x + b)
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return exp(x, ps[0], ps[1])
			},
			X: xdata,
			Y: ydata,
			N: 2,
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got, want := res.X, []float64{a, b}; !floats.EqualApprox(got, want, 0.1) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p := hplot.New()
		p.X.Label.Text = "exp(a*x+b)"
		p.Y.Label.Text = "y-data"
		p.Y.Min = 0
		p.Y.Max = 5
		p.X.Min = 0
		p.X.Max = 5

		s := hplot.NewS2D(hplot.ZipXY(xdata, ydata))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return exp(x, res.X[0], res.X[1])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err := p.Save(20*vg.Centimeter, -1, "testdata/exp-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_poly() {
	var (
		a    = 1.0
		b    = 2.0
		ps   = []float64{a, b}
		want = []float64{1.38592513, 1.98485122} // from scipy.curve_fit
	)

	poly := func(x float64, ps []float64) float64 {
		return ps[0] + ps[1]*x*x
	}

	xdata, ydata := genXY(100, poly, ps, -10, 10)

	res, err := fit.Curve1D(
		fit.Func1D{
			F:  poly,
			X:  xdata,
			Y:  ydata,
			Ps: []float64{1, 1},
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}

	if got := res.X; !floats.EqualApprox(got, want, 1e-6) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p := hplot.New()
		p.X.Label.Text = "f(x) = a + b*x*x"
		p.Y.Label.Text = "y-data"
		p.X.Min = -10
		p.X.Max = +10
		p.Y.Min = 0
		p.Y.Max = 220

		s := hplot.NewS2D(hplot.ZipXY(xdata, ydata))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return poly(x, res.X)
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err := p.Save(20*vg.Centimeter, -1, "testdata/poly-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_powerlaw() {
	var (
		amp   = 11.021171432949746
		index = -2.027389113217428
		want  = []float64{amp, index}
	)

	xdata, ydata, yerrs, err := readXYerr("testdata/powerlaw-data.txt")
	if err != nil {
		log.Fatal(err)
	}

	plaw := func(x, amp, index float64) float64 {
		return amp * math.Pow(x, index)
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return plaw(x, ps[0], ps[1])
			},
			X:   xdata,
			Y:   ydata,
			Err: yerrs,
			Ps:  []float64{1, 1},
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got := res.X; !floats.EqualApprox(got, want, 1e-3) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p := hplot.New()
		p.X.Label.Text = "f(x) = a * x^b"
		p.Y.Label.Text = "y-data"
		p.X.Min = 0
		p.X.Max = 10
		p.Y.Min = 0
		p.Y.Max = 10

		pts := make([]hbook.Point2D, len(xdata))
		for i := range pts {
			pts[i].X = xdata[i]
			pts[i].Y = ydata[i]
			pts[i].ErrY.Min = 0.5 * yerrs[i]
			pts[i].ErrY.Max = 0.5 * yerrs[i]
		}

		s := hplot.NewS2D(hbook.NewS2D(pts...), hplot.WithYErrBars(true))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return plaw(x, res.X[0], res.X[1])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err := p.Save(20*vg.Centimeter, -1, "testdata/powerlaw-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_hessian() {
	var (
		cst   = 3.0
		mean  = 30.0
		sigma = 20.0
		want  = []float64{cst, mean, sigma}
	)

	xdata, ydata, err := readXY("testdata/gauss-data.txt")
	if err != nil {
		log.Fatal(err)
	}

	// use a small sample
	xdata = xdata[:min(25, len(xdata))]
	ydata = ydata[:min(25, len(ydata))]

	gauss := func(x, cst, mu, sigma float64) float64 {
		v := (x - mu)
		return cst * math.Exp(-v*v/sigma)
	}

	f1d := fit.Func1D{
		F: func(x float64, ps []float64) float64 {
			return gauss(x, ps[0], ps[1], ps[2])
		},
		X:  xdata,
		Y:  ydata,
		Ps: []float64{10, 10, 10},
	}
	res, err := fit.Curve1D(f1d, nil, &optimize.NelderMead{})
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got := res.X; !floats.EqualApprox(got, want, 1e-3) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	inv := mat.NewSymDense(len(res.Location.X), nil)
	f1d.Hessian(inv, res.Location.X)
	// fmt.Printf("hessian: %1.2e\n", mat.Formatted(inv, mat.Prefix("         ")))

	popt := res.Location.X
	pcov := mat.NewDense(len(popt), len(popt), nil)
	{
		var chol mat.Cholesky
		if ok := chol.Factorize(inv); !ok {
			log.Fatalf("cov-matrix not positive semi-definite")
		}

		err := chol.InverseTo(inv)
		if err != nil {
			log.Fatalf("could not inverse matrix: %+v", err)
		}
		pcov.Copy(inv)
	}

	// compute goodness-of-fit.
	gof := newGoF(f1d.X, f1d.Y, popt, func(x float64) float64 {
		return f1d.F(x, popt)
	})

	pcov.Scale(gof.SSE/float64(len(f1d.X)-len(popt)), pcov)

	// fmt.Printf("pcov: %1.2e\n", mat.Formatted(pcov, mat.Prefix("      ")))

	var (
		n   = float64(len(f1d.X))    // number of data points
		ndf = n - float64(len(popt)) // number of degrees of freedom
		t   = distuv.StudentsT{
			Mu:    0,
			Sigma: 1,
			Nu:    ndf,
		}.Quantile(0.5 * (1 + 0.95))
	)

	for i, p := range popt {
		sigma := math.Sqrt(pcov.At(i, i))
		fmt.Printf("c%d: %1.5e [%1.5e, %1.5e] -- truth: %g\n", i, p, p-sigma*t, p+sigma*t, want[i])
	}
	// Output:
	//c0: 2.99999e+00 [2.99999e+00, 3.00000e+00] -- truth: 3
	//c1: 3.00000e+01 [3.00000e+01, 3.00000e+01] -- truth: 30
	//c2: 2.00000e+01 [2.00000e+01, 2.00000e+01] -- truth: 20
}

type GoF struct {
	SSE        float64 // Sum of squares due to error
	Rsquare    float64 // R-Square is the square of the correlation between the response values and the predicted response values
	NdF        int     // Number of degrees of freedom
	AdjRsquare float64 // Degrees of freedom adjusted R-Square
	RMSE       float64 // Root mean squared error
}

func newGoF(xs, ys, ps []float64, f func(float64) float64) GoF {
	switch {
	case len(xs) != len(ys):
		panic("invalid lengths")
	}

	var gof GoF

	var (
		ye = make([]float64, len(ys))
		nn = float64(len(xs) - 1)
		vv = float64(len(xs) - len(ps))
	)

	for i, x := range xs {
		ye[i] = f(x)
		dy := ys[i] - ye[i]
		gof.SSE += dy * dy
		gof.RMSE += dy * dy
	}

	gof.Rsquare = stat.RSquaredFrom(ye, ys, nil)
	gof.AdjRsquare = 1 - ((1 - gof.Rsquare) * nn / vv)
	gof.RMSE = math.Sqrt(gof.RMSE / float64(len(ys)-len(ps)))
	gof.NdF = len(ys) - len(ps)

	return gof
}
