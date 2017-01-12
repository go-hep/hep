# fit

`fit` is a WIP package to provide easy fitting models and curve fitting functions.

## Curve1D

### Fit a gaussian

![func1d-gaussian-example](https://github.com/go-hep/fit/raw/master/testdata/gauss-plot.png)
[embedmd]:# (fit_curve1d_test.go go /func ExampleCurve1D_gaussian/ /\n}/)
```go
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
```

### Fit an exponential

![func1d-curve-example](https://github.com/go-hep/fit/raw/master/testdata/curve-plot.png)
![func1d-ceres-example](https://github.com/go-hep/fit/raw/master/testdata/ceres-plot.png)
