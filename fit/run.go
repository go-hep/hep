// +build ignore

package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func fitFunc(x float64, p []float64) float64 {
	return p[0]*math.Cos(2*math.Pi/p[1]*x+p[2]) + p[3]*x
}

func readXY(fname string) (xs, ys []float64) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		xs = append(xs, x)

		y, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		ys = append(ys, y)
	}

	return
}

func main() {
	curveFit()
	ceresFit()
}

func curveFit() {
	xdata, ydata := readXY("testdata/curve-data.txt")
	log.Printf("fit(3): %v\n", efit(3, 2.5, 1.3))
	fmt.Printf("--- LLS ---\n")
	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return efit(x, ps[0], ps[1])
			},
			N: 2,
			X: xdata,
			Y: ydata,
		},
		nil, &optimize.LBFGS{},
	)
	switch err {
	case nil:
		fmt.Printf("res.X: %v\n", res.X)
		fmt.Printf("res.F: %v\n", res.F)
		fmt.Printf("res.G: %v\n", res.Gradient)
		fmt.Printf("res.H: %v\n", res.Hessian)
		fmt.Printf("res.Stats: %+v\n", res.Stats)
		fmt.Printf("fit(3): %v\n", efit(3, res.X[0], res.X[1]))
	default:
		fmt.Printf("err: %v\n", err)
	}
	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "x"
		p.Y.Label.Text = "curve data"
		p.Add(plotter.NewGrid())
		s, err := plotter.NewScatter(hplot.ZipXY(xdata, ydata))
		if err != nil {
			log.Fatal(err)
		}
		s.Color = color.RGBA{255, 0, 0, 255}

		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return efit(x, res.X[0], res.X[1])
		})
		f.Color = color.RGBA{0, 0, 255, 255}
		p.Add(f)

		p.Save(20*vg.Centimeter, 15*vg.Centimeter, "testdata/curve-plot.png")
	}

}

func efit(x, a, b float64) float64 {
	return a*math.Exp(-x) + b
}

func ceresFitFunc(x, a, b float64) float64 {
	return math.Exp(a*x + b)
}

func ceresFit() {
	xdata, ydata := readXY("testdata/ceres-data.txt")
	fmt.Printf("--- Ceres FIT ---\n")
	res, err := fit.Curve1D(
		fit.Func1D{
			F:  func(x float64, ps []float64) float64 { return ceresFitFunc(x, ps[0], ps[1]) },
			X:  xdata,
			Y:  ydata,
			Ps: []float64{1, 1},
		}, nil, nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("res.Status: %v\n", res.Status)
	fmt.Printf("res.X: %v\n", res.X)
	fmt.Printf("res.F: %v\n", res.F)
	fmt.Printf("res.G: %v\n", res.Gradient)
	fmt.Printf("res.H: %v\n", res.Hessian)
	fmt.Printf("res.Stats: %+v\n", res.Stats)

	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "x"
		p.Y.Label.Text = "ceres data"
		p.Add(plotter.NewGrid())
		s, err := plotter.NewScatter(hplot.ZipXY(xdata, ydata))
		if err != nil {
			log.Fatal(err)
		}
		s.Color = color.RGBA{255, 0, 0, 255}

		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return ceresFitFunc(x, res.X[0], res.X[1])
		})
		f.Color = color.RGBA{0, 0, 255, 255}
		p.Add(f)

		p.Save(20*vg.Centimeter, 15*vg.Centimeter, "testdata/ceres-plot.png")
	}

}
