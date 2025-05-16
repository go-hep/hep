# fit

[![GoDoc](https://godoc.org/go-hep.org/x/hep/fit?status.svg)](https://godoc.org/go-hep.org/x/hep/fit)

`fit` is a WIP package to provide easy fitting models and curve fitting functions.

## H1D

### Fit a gaussian

![h1d-gaussian-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/h1d-gauss-plot_golden.png)

[embedmd]:# (hist_example_test.go go /func ExampleH1D_gaussian/ /\n}/)
```go
func ExampleH1D_gaussian() {
	var (
		mean  = 2.0
		sigma = 4.0
		// Values from gonum/optimize:
		want = []float64{447.0483517081991, 2.02127773281178, 3.9965893891862687}
		// Values from ROOT:
		// want  = []float64{4.53720e+02, 1.93218e+00, 3.93188e+00}
	)

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    mean,
		Sigma: sigma,
		Src:   rand.New(rand.NewSource(0)),
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
		p.X.Label.Text = "f(x) = cst * exp(-0.5 * ((x-mu)/sigma)^2)"
		p.Y.Label.Text = "y-data"
		p.Y.Min = 0

		h := hplot.NewH1D(hist)
		h.Color = color.RGBA{0, 0, 255, 255}
		p.Add(h)

		f := plotter.NewFunction(func(x float64) float64 {
			return gauss(x, res.X[0], res.X[1], res.X[2])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err := p.Save(20*vg.Centimeter, -1, "testdata/h1d-gauss-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

## Curve1D

### Fit a gaussian

![func1d-gaussian-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/gauss-plot_golden.png)

[embedmd]:# (curve1d_example_test.go go /func ExampleCurve1D_gaussian/ /\n}/)
```go
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
```

### Fit a powerlaw (with Y-errors)

![func1d-powerlaw-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/powerlaw-plot_golden.png)

[embedmd]:# (curve1d_example_test.go go /func ExampleCurve1D_powerlaw/ /\n}/)
```go
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
```

### Fit an exponential

![func1d-exp-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/exp-plot_golden.png)

[embedmd]:# (curve1d_example_test.go go /func ExampleCurve1D_exponential/ /\n}/)
```go
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
```

### Fit a polynomial

![func1d-poly-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/poly-plot_golden.png)

[embedmd]:# (curve1d_example_test.go go /func ExampleCurve1D_poly/ /\n}/)
```go
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
```
## Fitting with more than one independent variable (x has more than one dimension)
### Fit a flat plane

![2d-example](https://codeberg.org/go-hep/hep/raw/branch/main/fit/testdata/2d-plane-plot_golden.png)

[embedmd]:# (curve_nd_example_test.go go /func ExampleCurveND_plane/ /\n}/)
```go
func ExampleCurveND_plane() {
	var (
		m1    = 0.3
		m2    = 0.1
		c     = 0.2
		ps    = []float64{m1, m2, c}
		n0    = 10
		n1    = 10
		x0min = -1.
		x0max = 1.
		x1min = -1.
		x1max = 1.
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
		// slicing for a particular x0 value to plot y as a function of x1,
		// to visualise how well the fit is working for a given x0.
		x0Selection := 8
		if 0 > x0Selection || x0Selection > n0 {
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

		p := hplot.New()
		p.Title.Text = fmt.Sprintf("Slice of plane at x0 = %.2f", x0SlicePos)
		p.X.Label.Text = "x1"
		p.Y.Label.Text = "y"
		p.Y.Min = x1min
		p.Y.Max = x1max
		p.X.Min = x0min
		p.X.Max = x0max

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
		err := p.Save(20*vg.Centimeter, -1, "testdata/2d-plane-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}
```
